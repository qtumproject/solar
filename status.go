package solar

import (
	"fmt"
	"os"
)

func init() {
	_ = app.Command("status", "Show statuses of contracts")

	appTasks["status"] = func() (err error) {
		repo := solar.ContractsRepository()

		if len(repo.contracts) == 0 {
			fmt.Println("No deployed contract yet")
			os.Exit(0)
		}

		for _, contract := range repo.SortedContracts() {
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
