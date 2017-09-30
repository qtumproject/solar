package solar

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

func init() {
	cmd := app.Command("encode", "ABI encoding for a method call")

	contractName := cmd.Arg("contractName", "Name of contract").Required().String()
	methodName := cmd.Arg("methodName", "Method name of contract").Required().String()
	jsonParams := cmd.Arg("jsonParams", "Parameters as a json array").Default("[]").String()

	appTasks["encode"] = func() (err error) {
		repo := solar.ContractsRepository()

		c, ok := repo.contracts[*contractName]
		if !ok {
			return errors.Errorf("Cannot find contract: %s", *contractName)
		}

		abi, err := c.encodingABI()
		if err != nil {
			return
		}

		var params []interface{}
		if jsonParams != nil {
			err := json.Unmarshal([]byte(*jsonParams), &params)
			if err != nil {
				return err
			}
		}

		data, err := abi.Pack(*methodName, params...)
		if err != nil {
			return err
		}

		fmt.Println(hex.EncodeToString(data))

		return nil
	}
}
