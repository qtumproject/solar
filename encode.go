package solar

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/qtumproject/solar/contract"
)

func init() {
	cmd := app.Command("encode", "ABI encoding for a method call")

	contractName := cmd.Arg("contractName", "Name of contract").Required().String()
	methodName := cmd.Arg("methodName", "Method name of contract").Required().String()
	jsonParams := cmd.Arg("jsonParams", "Parameters as a json array").Default("[]").String()

	appTasks["encode"] = func() (err error) {
		repo := solar.ContractsRepository()

		c, ok := repo.Get(*contractName)
		if !ok {
			return errors.Errorf("Cannot find contract: %s", *contractName)
		}

		abi, err := c.EncodingABI()
		if err != nil {
			return
		}

		var params []interface{}
		if jsonParams != nil {
			jsonParams := solar.ExpandJSONParams(*jsonParams)

			err := json.Unmarshal([]byte(jsonParams), &params)
			if err != nil {
				return err
			}
		}

		data, err := abi.Pack(*methodName, params...)
		if err != nil {
			return err
		}

		fmt.Println(contract.Bytes(data).String())

		return nil
	}
}
