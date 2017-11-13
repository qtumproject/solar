package solar

import (
	"fmt"
	"log"
)

func init() {
	cmd := app.Command("deploy", "Compile Solidity contracts.")

	force := cmd.Flag("force", "Overwrite previously deployed contract with the same name").Bool()
	noconfirm := cmd.Flag("no-confirm", "Don't wait for network to confirm deploy").Bool()
	noFastConfirm := cmd.Flag("no-fast-confirm", "(dev) Don't generate block to confirm deploy immediately").Bool()

	sourceFilePath := cmd.Arg("source", "Solidity contracts to deploy.").Required().String()
	name := cmd.Arg("name", "Name of contract").Required().String()
	jsonParams := cmd.Arg("jsonParams", "Constructor params as a json array").Default("").String()

	appTasks["deploy"] = func() (err error) {
		deployer := solar.Deployer()

		fmt.Printf("   \033[36mdeploy\033[0m %s => %s\n", *sourceFilePath, *name)

		err = deployer.CreateContract(*name, *sourceFilePath, []byte(*jsonParams), *force)
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

type Deployer struct {
	rpc  *qtumRPC
	repo *contractsRepository
}

func (d *Deployer) CreateContract(name, filepath string, jsonParams []byte, overwrite bool) (err error) {
	// if !overwrite && d.repo.Exists(name) {
	// 	return errors.Errorf("name already used: %s", name)
	// }

	// gasLimit := 300000

	// rpc := d.rpc

	// contract, err := compileSource(filepath, CompilerOptions{})
	// if err != nil {
	// 	return errors.Wrap(err, "compile")
	// }

	// abi, err := contract.encodingABI()
	// if err != nil {
	// 	return errors.Wrap(err, "abi")
	// }

	// constructor := abi.Constructor

	// if len(constructor.Inputs) == 0 && len(jsonParams) != 0 {
	// 	return errors.New("does not expect constructor params")
	// }

	// bin := contract.Bin
	// if len(constructor.Inputs) != 0 {
	// 	var params []interface{}
	// 	err = json.Unmarshal(jsonParams, &params)
	// 	if err != nil {
	// 		return errors.Errorf("expected constructor params in JSON, got: %#v", string(jsonParams))
	// 	}

	// 	packedParams, err := abi.Constructor.Pack(params...)
	// 	if err != nil {
	// 		return errors.Wrap(err, "constructor")
	// 	}

	// 	bin = append(bin, packedParams...)
	// }

	// var tx TransactionReceipt

	// err = rpc.Call(&tx, "createcontract", bin, gasLimit)

	// if err != nil {
	// 	return errors.Wrap(err, "createcontract")
	// }

	// // fmt.Println("tx", tx.Address)
	// // fmt.Println("contract name", contract.Name)

	// deployedContract := &DeployedContract{
	// 	Name:             contract.Name,
	// 	DeployName:       name,
	// 	CompiledContract: *contract,
	// 	TransactionID:    tx.TxID,
	// 	Address:          tx.Address,
	// 	CreatedAt:        time.Now(),
	// }

	// err = d.repo.Set(name, deployedContract)
	// if err != nil {
	// 	return
	// }

	// err = d.repo.Commit()
	// if err != nil {
	// 	return
	// }
	// pretty.Println("rpc err", string(res.RawError))

	return nil
}
