package solar

import (
	"fmt"
	"os"
)

func init() {
	cli := app.Command("status", "Show statuses of contracts")
	contractNames := cli.Arg("names", "contract names").Strings()

	appTasks["status"] = func() (err error) {
		names := *contractNames
		repo := solar.ContractsRepository()

		if len(repo.contracts) == 0 {
			fmt.Println("No deployed contract yet")
			os.Exit(0)
		}

		var contracts []*DeployedContract
		if len(names) != 0 {
			for _, name := range names {
				contract, found := repo.contracts[name]
				if !found {
					fmt.Printf("\u2757\ufe0f %s: not found\n", name)
					continue
				}
				contracts = append(contracts, contract)
			}
			// contracts =
		} else {
			contracts = repo.SortedContracts()
		}

		for _, contract := range contracts {
			// FIXME: store deploy name in contract
			name := contract.DeployName
			if contract.Confirmed {
				fmt.Printf("\u2705  %s\n", name)
			} else {
				fmt.Printf("   %s\n", name)
			}

			fmt.Printf("        txid: %s\n", contract.TransactionID)
			fmt.Printf("     address: %s\n", contract.Address)
			fmt.Printf("   confirmed: %v\n", contract.Confirmed)

			fmt.Println("")

		}

		return nil
	}
}
