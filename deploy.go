package solar

import (
	"fmt"
	"log"
	"strings"

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

		contract, err := compiler.Compile()
		if err != nil {
			return errors.Wrap(err, "compile")
		}

		deployer := solar.Deployer()

		fmt.Printf("   \033[36mdeploy\033[0m %s => %s\n", target.file, target.name)

		err = deployer.CreateContract(contract, []byte(*jsonParams), target.name, *force, *aslib)
		if err != nil {
			fmt.Println("\u2757\ufe0f \033[36mdeploy\033[0m", err)
			return
		}

		newContracts := repo.UnconfirmedContracts()
		if *noconfirm == false && len(newContracts) != 0 {
			// Force local chain to generate a block immediately.
			allowFastConfirm := *solarEnv == "development" || *solarEnv == "test"
			if *noFastConfirm == false && allowFastConfirm {
				//fmt.Println("call deployer.Mine")
				err = deployer.Mine()
				if err != nil {
					log.Println(err)
				}
			}

			err := repo.ConfirmAll(getConfirmUpdateProgressFunc(), deployer.ConfirmContract)
			if err != nil {
				return err
			}
		}

		return
	}
}
