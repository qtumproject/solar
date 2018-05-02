package solar

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/qtumproject/solar/deployer"

	"github.com/qtumproject/solar/contract"

	"github.com/pkg/errors"
)

type deployTarget struct {
	file string
	name string
}

func parseDeployTarget(target string) deployTarget {
	parts := strings.Split(target, ":")

	if len(parts) == 1 {
		return deployTarget{
			file: parts[0],
			name: parts[0],
		}
	}

	// deploy name by default is the file path relative to project root
	return deployTarget{
		file: parts[0],
		name: parts[1],
	}
}

func init() {
	cmd := app.Command("deploy", "Compile Solidity contracts.")

	force := cmd.Flag("force", "Overwrite previously deployed contract with the same deploy name").Bool()
	aslib := cmd.Flag("lib", "Deploy the contract as a library").Bool()
	noconfirm := cmd.Flag("no-confirm", "Don't wait for network to confirm deploy").Bool()
	noFastConfirm := cmd.Flag("no-fast-confirm", "(dev) Don't generate block to confirm deploy immediately").Bool()
	gasLimit := cmd.Flag("gasLimit", "gas limit for creating a contract").Default("3000000").Uint()
	gasPrice := cmd.Flag("gasPrice", "gas price for transaction in satoshi (default 40) or gwei (default 1)").Default("").String()

	target := cmd.Arg("target", "Solidity contracts to deploy.").Required().String()
	jsonParams := cmd.Arg("jsonParams", "Constructor params as a json array").Default("").String()

	appTasks["deploy"] = func() (err error) {
		target := parseDeployTarget(*target)

		opts, err := solar.SolcOptions()
		if err != nil {
			return
		}

		filename := target.file

		repo := solar.ContractsRepository()

		compiler := Compiler{
			Opts:     *opts,
			Filename: filename,
			Repo:     repo,
		}

		compiledContract, err := compiler.Compile()
		if err != nil {
			return errors.Wrap(err, "compile")
		}

		var params []byte
		if jsonParams != nil {
			jsonParams := solar.ExpandJSONParams(*jsonParams)

			params = []byte(jsonParams)
		}

		gasPrice := *gasPrice
		if gasPrice == "" {
			switch solar.RPCPlatform() {
			case RPCEthereum:
				gasPrice = "1"
			case RPCQtum:
				gasPrice = "40"
			}
		}

		parsedGasPrice, _, err := big.ParseFloat(gasPrice, 0, big.MaxPrec, big.ToNearestEven)
		if err != nil {
			return errors.Errorf("Cannot parse gas price: %s", gasPrice)
		}

		if parsedGasPrice.Sign() <= 0 {
			return errors.Errorf("Gas price must be positive: %s", gasPrice)
		}

		fmt.Println("cli gasPrice", gasPrice, parsedGasPrice.String())

		deployOpts := deployer.Options{
			Name:      target.name,
			Overwrite: *force,
			AsLib:     *aslib,
			GasLimit:  *gasLimit,
			GasPrice:  parsedGasPrice,
		}

		dpl := solar.Deployer()

		err = dpl.CreateContract(compiledContract, params, &deployOpts)
		if err != nil {
			fmt.Println("\u2757\ufe0f \033[36mdeploy\033[0m", err)
			return
		}

		// Add related contracts to repo
		relatedContracts, err := compiler.RelatedContracts()
		if err != nil {
			return err
		}

		if len(relatedContracts) > 0 {
			for name, c := range relatedContracts {
				repo.Related[name] = c
			}

			err = repo.Commit()
			if err != nil {
				return
			}
		}

		newContracts := repo.UnconfirmedContracts()
		if *noconfirm == false && len(newContracts) != 0 {
			// Force local chain to generate a block immediately.
			allowFastConfirm := *solarEnv == "development" || *solarEnv == "test"
			if *noFastConfirm == false && allowFastConfirm {
				//fmt.Println("call deployer.Mine")
				err = dpl.Mine()
				if err != nil {
					log.Println(err)
				}
			}

			err := repo.ConfirmAll(getConfirmUpdateProgressFunc(), dpl.ConfirmContract)
			if err != nil {
				return err
			}

			var deployedContract *contract.DeployedContract
			if *aslib {
				deployedContract, _ = repo.GetLib(target.name)
			} else {
				deployedContract, _ = repo.Get(target.name)
			}

			if deployedContract == nil {
				return errors.New("failed to deploy contract")
			}

			fmt.Printf("   \033[36mdeployed\033[0m %s => %s\n", target.name, deployedContract.Address)
		}

		return
	}
}
