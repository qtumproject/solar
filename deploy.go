package solar

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

func init() {
	cmd := app.Command("deploy", "Compile Solidity contracts.")

	force := cmd.Flag("force", "Overwrite previously deployed contract with the same name").Bool()
	noconfirm := cmd.Flag("no-confirm", "Don't wait for network to confirm deploy").Bool()
	noFastConfirm := cmd.Flag("no-fast-confirm", "(dev) Don't generate block to confirm deploy immediately").Bool()

	sourceFilePath := cmd.Arg("source", "Solidity contracts to deploy.").Required().String()
	name := cmd.Arg("name", "Name of contract").Required().String()
	// jsonParams := cmd.Arg("jsonParams", "Parameters as a json array").Default("[]").String()

	appTasks["deploy"] = func() (err error) {
		deployer := solar.Deployer()

		fmt.Printf("   \033[36mdeploy\033[0m %s => %s\n", *sourceFilePath, *name)

		err = deployer.CreateContract(*name, *sourceFilePath, *force)
		if err != nil {
			fmt.Println("\u2757\ufe0f \033[36mdeploy\033[0m", err)
			return
		}

		repo := solar.ContractsRepository()
		newContracts := repo.UnconfirmedContracts()

		if *noconfirm == false && len(newContracts) != 0 {
			// Force local chain to generate a block immediately.
			allowFastConfirm := *solarEnv == "development" || *solarEnv == "test"
			if *noFastConfirm == false && allowFastConfirm {
				rpc := solar.RPC()
				err := rpc.Call(nil, "generate", 1)
				if err != nil {
					log.Println(err)
				}
			}

			err := repo.ConfirmAll()
			if err != nil {
				return err
			}
		}

		return
	}
}

// type deployTarget struct {
// 	Name     string
// 	FilePath string
// }

// func parseDeployTarget(target string) deployTarget {
// 	parts := strings.Split(target, ":")

// 	filepath := parts[0]

// 	var name string
// 	if len(parts) == 2 {
// 		name = parts[1]
// 	} else {
// 		name = stringLowerFirstRune(basenameNoExt(filepath))
// 	}

// 	// TODO verify name for valid JS name

// 	t := deployTarget{
// 		Name:     name,
// 		FilePath: filepath,
// 	}

// 	return t
// }

type Deployer struct {
	rpc  *qtumRPC
	repo *contractsRepository
}

func (d *Deployer) CreateContract(name, filepath string, overwrite bool) (err error) {
	if !overwrite && d.repo.Exists(name) {
		return errors.Errorf("name already used: %s", name)
	}

	gasLimit := 300000

	rpc := d.rpc

	contract, err := compileSource(filepath, CompilerOptions{})
	if err != nil {
		return errors.Wrap(err, "compile")
	}

	var tx TransactionReceipt

	err = rpc.Call(&tx, "createcontract", contract.Bin.String(), gasLimit)

	if err != nil {
		return errors.Wrap(err, "createcontract")
	}

	// fmt.Println("tx", tx.Address)
	// fmt.Println("contract name", contract.Name)

	deployedContract := &DeployedContract{
		Name:             contract.Name,
		DeployName:       name,
		CompiledContract: *contract,
		TransactionID:    tx.TxID,
		Address:          tx.Address,
		CreatedAt:        time.Now(),
	}

	err = d.repo.Set(name, deployedContract)
	if err != nil {
		return
	}

	err = d.repo.Commit()
	if err != nil {
		return
	}
	// pretty.Println("rpc err", string(res.RawError))

	return nil
}
