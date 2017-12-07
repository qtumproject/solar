package solar

import (
	"fmt"
	"os"

	"github.com/qtumproject/solar/contract"
)

func init() {
	cli := app.Command("status", "Show statuses of contracts")
	contractNames := cli.Arg("names", "contract names").Strings()

	appTasks["status"] = func() (err error) {
		names := *contractNames
		repo := solar.ContractsRepository()

		if repo.Len() == 0 {
			fmt.Println("No deployed contract yet")
			os.Exit(0)
		}

		var contracts []*contract.DeployedContract
		if len(names) != 0 {
			for _, name := range names {
				c, found := repo.Get(name)
				if !found {
					fmt.Printf("\u2757\ufe0f %s: not found\n", name)
					continue
				}
				contracts = append(contracts, c)
			}
			// contracts =
		} else {
			contracts = repo.SortedContracts()
		}

		for _, c := range contracts {
			// FIXME: store deploy name in contract
			name := c.DeployName
			if c.Confirmed {
				fmt.Printf("\u2705  %s\n", name)
			} else {
				fmt.Printf("   %s\n", name)
			}

			fmt.Printf("        txid: %s\n", c.TransactionID)
			fmt.Printf("     address: %s\n", c.Address)
			fmt.Printf("   confirmed: %v\n", c.Confirmed)

			fmt.Println("")

		}

		return nil
	}
}
