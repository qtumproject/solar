package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/crypto/sha3"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("solar", "Solidity smart contract deployment management.")
	buildCmd    = app.Command("build", "Compile Solidity contracts.")
	buildTarget = buildCmd.Arg("target", "Source file or directory. Default is the current directory (`.`)").Default(".").String()
)

type Builder struct {
	compilerOpts CompilerOptions
	outputDir    string
}

func (b *Builder) Compile(filename string) error {
	outputFilename := filename + ".json"

	compiledContracts, err := compileSourceFile(filename, b.compilerOpts)
	if err != nil {
		return errors.Wrap(err, "compile")
	}

	outf, err := os.Create(outputFilename)
	if err != nil {
		return errors.Wrap(err, "output")
	}
	defer outf.Close()

	contracts := make(map[string]CompiledContract)
	for _, contract := range compiledContracts {
		contracts[contract.Name] = contract
	}

	enc := json.NewEncoder(outf)
	enc.SetIndent("", "\t")
	err = enc.Encode(contracts)
	if err != nil {
		return err
	}

	return nil
}

func buildTask() (err error) {
	fi, err := os.Stat(*buildTarget)

	if os.IsNotExist(err) {
		return errors.Errorf("Compile target not found: %s", *buildTarget)
	}

	if err != nil {
		return
	}

	builder := Builder{
		compilerOpts: CompilerOptions{},
	}

	if fi.IsDir() {
		pat := path.Join(*buildTarget, "*.sol")
		matches, err := filepath.Glob(pat)
		if err != nil {
			return errors.Wrap(err, "glob")
		}

		for _, filename := range matches {
			fmt.Println("Compiling:", filename)
			err := builder.Compile(filename)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		err := builder.Compile(*buildTarget)
		if err != nil {
			fmt.Println(err)
		}
	}

	return
}

var cliTasks = map[string]func() error{
	"build": buildTask,
}

func main() {
	cmdName, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	task := cliTasks[cmdName]
	err = task()
	if err != nil {
		log.Fatalln(err)
	}
}

type ABIDefinition struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Payable  bool      `json:"payable"`
	Inputs   []ABIType `json:"inputs"`
	Outputs  []ABIType `json:"outputs"`
	Constant bool      `json:"constant"`
}

type ABIType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Bytes []byte

func (b Bytes) String() string {
	return hex.EncodeToString(b)
}
func (b Bytes) MarshalJSON() ([]byte, error) {
	hexstr := fmt.Sprintf("\"%s\"", b.String())
	return []byte(hexstr), nil
}

type CompilerOptions struct {
	NoOptimize bool
}

type CompiledContract struct {
	Name string          `json:"name"`
	ABI  []ABIDefinition `json:"abi"`
	Bin  Bytes           `json:"bin"`
	// KECAAK256 of bytecode without auxdata
	BinKeccak256 Bytes `json:"binhash"`
}

type rawCompilerOutput struct {
	Version   string
	Contracts map[string]rawCompiledContract
}

type rawCompiledContract struct {
	RawMetadata string `json:"metadata"`
	Bin         []byte
	Metadata    struct {
		Output struct {
			Version  int64
			Language string
			ABI      []ABIDefinition `json:"abi"`
		}
	}
}

func (c *rawCompiledContract) BinHash256() []byte {
	bin := c.BinWithoutAuxData()
	h := sha3.NewKeccak256()
	h.Write(bin)
	binDigest := h.Sum(nil)
	return binDigest
}

func (c *rawCompiledContract) BinWithoutAuxData() []byte {
	// https://solidity.readthedocs.io/en/develop/miscellaneous.html#encoding-of-the-metadata-hash-in-the-bytecode
	// 0xa1 0x65 'b' 'z' 'z' 'r' '0' 0x58 0x20 <32 bytes swarm hash> 0x00 0x29
	// a1 65 62 7a 7a 72 30 58 20 [32 bytes] 0x00 0x29
	// 11 + 32 bytes

	return c.Bin[0 : len(c.Bin)-11-32-1]
}

func (c *rawCompiledContract) UnmarshalJSON(data []byte) error {
	type dataStruct struct {
		RawMetadata string `json:"metadata"`
		BinStr      string `json:"bin"`
	}

	var dest dataStruct

	err := json.Unmarshal(data, &dest)
	if err != nil {
		return err
	}

	// Recursively parse Metadata, which is a json string.
	err = json.Unmarshal([]byte(dest.RawMetadata), &c.Metadata)
	if err != nil {
		return err
	}

	bin, err := hex.DecodeString(dest.BinStr)
	if err != nil {
		return err
	}

	c.RawMetadata = dest.RawMetadata
	c.Bin = bin

	return nil
}

func compileSourceFile(src string, opts CompilerOptions) (compiledContracts []CompiledContract, err error) {
	f, err := os.Open(src)
	if err != nil {
		return
	}
	defer f.Close()

	args := []string{"-", "--combined", "bin,metadata"}

	if !opts.NoOptimize {
		args = append(args, "--optimize")
	}

	// fmt.Printf("exec: solc %v\n", args)
	cmd := exec.Command("solc", args...)
	cmd.Stdin = f
	output, err := cmd.Output()
	if err != nil {
		return
	}

	// log.Println("output", string(output))

	var compilerOutput rawCompilerOutput
	err = json.Unmarshal(output, &compilerOutput)
	if err != nil {
		return nil, errors.Wrap(err, "parse output")
	}

	// fmt.Printf("%#v", compilerOutput)
	for name, c := range compilerOutput.Contracts {
		// fmt.Println(name, c.RawMetadata)

		// name: <stdin>:ContractName
		contractName := name
		parts := strings.Split(name, ":")
		if len(parts) == 2 {
			contractName = parts[1]
		}

		// log.Println("bin", c.Bin)

		compiledContract := CompiledContract{
			Name:         contractName,
			Bin:          c.Bin,
			BinKeccak256: c.BinHash256(),
			ABI:          c.Metadata.Output.ABI,
		}

		// pretty.Println("abi", c.Metadata.Output.ABI)
		// fmt.Println(cc)
		compiledContracts = append(compiledContracts, compiledContract)
	}

	return
}
