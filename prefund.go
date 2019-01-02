package solar

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	cli := app.Command("prefund", "(qtum) fund an owner address with utxos")

	owner := cli.Arg("owner", "contract name or address to fund").Required().String()

	amount := cli.Arg("amount", "fund an utxo with this amount").Required().Float64()
	multiples := cli.Arg("multiples", "fund this number of identical utxos").Default("1").Int()

	appTasks["prefund"] = func() (err error) {
		rpc := solar.QtumRPC()

		repo := solar.ContractsRepository()

		var ownerAddr string

		contract, found := repo.Get(*owner)
		if found {
			ownerAddr = contract.Sender
		} else {
			ownerAddr = *owner
		}

		// if the address is hexadecimal, convert it to base58 address
		_, err = hex.DecodeString(ownerAddr)
		if err == nil {
			var b58addr string
			rpcErr := rpc.Call(&b58addr, "fromhexaddress", ownerAddr)
			if rpcErr != nil {
				return errors.Wrap(err, "convert hex address")
			}

			ownerAddr = b58addr
		}

		// The JSON object is allowed to have duplicate keys for this call
		// { <addr>: <amount>, ... }

		var utxos []string
		for i := 0; i < *multiples; i++ {
			utxo := fmt.Sprintf(`"%s": %f`, ownerAddr, *amount)
			utxos = append(utxos, utxo)
		}

		amounts := "{\n" + strings.Join(utxos, ",\n") + "\n}"

		// fmt.Println("json utxos", amounts)

		var result interface{}

		/*
			sendmanywithdupes "fromaccount" {"address":amount,...} ( minconf "comment" ["address",...] )

			1. "fromaccount"         (string, required) DEPRECATED. The account to send the funds from. Should be "" for the default account
			2. "amounts"             (string, required) A json object with addresses and amounts
			    {
			      "address":amount   (numeric or string) The qtum address is the key, the numeric amount (can be string) in QTUM is the value
			      ,...
			    }
		*/
		err = rpc.Call(&result, "sendmanywithdupes", "", json.RawMessage(amounts))
		if err != nil {
			return
		}

		fmt.Println("sendmanywithdupes txid:", result)

		return
	}
}
