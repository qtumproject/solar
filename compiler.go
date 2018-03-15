package solar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/qtumproject/solar/contract"
)

type rawCompilerOutput struct {
	Version   string
	Contracts map[string]contract.RawCompiledContract
}

func (o *rawCompilerOutput) CompiledContracts() contract.CompiledContracts {
	contracts := make(contract.CompiledContracts)

	for name, rawContract := range o.Contracts {
		if len(rawContract.Bin) == 0 {
			continue
		}

		// name: filepath:ContractName
		contractName := name
		parts := strings.Split(name, ":")
		if len(parts) == 2 {
			contractName = parts[1]
		}

		compiledContract := &contract.CompiledContract{
			Source:       name,
			Name:         contractName,
			Bin:          rawContract.Bin,
			BinKeccak256: rawContract.BinHash256(),
			ABI:          rawContract.Metadata.Output.ABI,
		}

		contracts[name] = compiledContract
	}

	return contracts
}

type CompilerError struct {
	SourceFile  string
	ErrorOutput string
}

func (err *CompilerError) Error() string {
	return err.ErrorOutput
}

type CompilerOptions struct {
	NoOptimize bool
	AllowPaths []string
}

type Compiler struct {
	// only used for error reporting
	Filename string
	Opts     CompilerOptions
	Repo     *contract.ContractsRepository

	compiledContracts contract.CompiledContracts
}

func (c *Compiler) mainContractSource() string {
	mainContractName := basenameNoExt(c.Filename)

	return c.Filename + ":" + mainContractName
}

func (c *Compiler) mainContractName() string {
	mainContractName := basenameNoExt(c.Filename)
	return mainContractName
}

// Compile returns only the contract that has the same name as the source file
func (c *Compiler) Compile() (*contract.CompiledContract, error) {

	contracts, err := c.getCompiledContracts()
	if err != nil {
		return nil, err
	}

	contractSrc := c.mainContractSource()
	contract, ok := contracts[contractSrc]
	if !ok {
		return nil, errors.Errorf("cannot find contract: %s", c.mainContractName())
	}

	return contract, nil
}

func (c *Compiler) RelatedContracts() (contract.CompiledContracts, error) {
	contracts, err := c.getCompiledContracts()
	if err != nil {
		return nil, err
	}

	contractSrc := c.mainContractSource()
	relatedContracts := make(contract.CompiledContracts)
	for name, contract := range contracts {
		if name == contractSrc {
			continue
		}

		relatedContracts[name] = contract
	}

	return relatedContracts, nil
}

// CompileAll returns all contracts in a source file
func (c *Compiler) compileAll() (contract.CompiledContracts, error) {
	_, err := os.Stat(c.Filename)

	if err != nil && os.IsNotExist(err) {
		return nil, errors.Errorf("file not found: %s", c.Filename)
	}

	output, err := c.execSolc()
	if err != nil {
		return nil, err
	}

	return output.CompiledContracts(), nil
}

func (c *Compiler) getCompiledContracts() (contract.CompiledContracts, error) {
	if c.compiledContracts != nil {
		return c.compiledContracts, nil
	}

	contracts, err := c.compileAll()
	if err != nil {
		return nil, err
	}

	c.compiledContracts = contracts
	return contracts, nil
}

func (c *Compiler) execSolc() (*rawCompilerOutput, error) {
	opts := c.Opts

	filename := c.Filename

	args := []string{filename, "--combined-json", "bin,metadata"}

	if !opts.NoOptimize {
		args = append(args, "--optimize")
	}

	if len(opts.AllowPaths) > 0 {
		args = append(args, "--allow-paths", strings.Join(opts.AllowPaths, ","))
	}

	// libraries linkage support
	if c.Repo != nil && len(c.Repo.Libraries) > 0 {
		var linkages []string
		// A linkable library is specified with a string that looks like:
		//    contracts/SafeMathLib.sol:SafeMathLib:4242424242424242424242424242424242424242
		for _, lib := range c.Repo.Libraries {
			linkages = append(linkages, fmt.Sprintf("%s:%s:%s", lib.DeployName, lib.Name, lib.Address))
		}

		args = append(args, "--libraries", strings.Join(linkages, ","))
	}

	var stderr bytes.Buffer

	fmt.Printf("exec: solc %v\n", args)
	cmd := exec.Command("solc", args...)
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()

	if _, hasExitErr := err.(*exec.ExitError); hasExitErr {
		return nil, &CompilerError{
			SourceFile:  filename,
			ErrorOutput: stderr.String(),
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "exec")
	}

	output := &rawCompilerOutput{}
	// fmt.Println("solc output", string(stdout))
	err = json.Unmarshal(stdout, output)
	if err != nil {
		return nil, errors.Wrap(err, "parse solc output")
	}

	return output, nil
}
