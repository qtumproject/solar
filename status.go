package solar

import (
	"fmt"
	"os"
)

func init() {
	_ = app.Command("status", "Compile Solidity contracts.")

	appTasks["status"] = func() (err error) {
		repo := solar.ContractsRepository()
		rpc := solar.RPC()

		if len(repo.contracts) == 0 {
			fmt.Println("No deployed contract yet")
			os.Exit(0)
		}

		for name, contract := range repo.contracts {
			if contract.Confirmed {
				fmt.Printf("%s\t%s\tconfirmed\n", name, contract.Address)
				continue
			}

			result := make(map[string]interface{})
			err := rpc.Call(&result, "getaccountinfo", contract.Address)
			if err != nil {
				fmt.Printf("%s\t%s\n", name, err)
				continue
			}

			fmt.Printf("%s\t%s\tconfirmed\n", name, contract.Address)
		}

		return nil
	}
}
