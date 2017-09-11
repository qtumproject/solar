package solar

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

func init() {
	_ = app.Command("status", "Compile Solidity contracts.")

	appTasks["status"] = func() (err error) {
		repo, err := openDeployedContractsRepository("solar.json")
		if err != nil {
			return
		}

		rpcURL, err := url.Parse(*solarRPC)
		if err != nil {
			return errors.Wrap(err, "rpc host")
		}

		rpc := qtumRPC{rpcURL}

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
