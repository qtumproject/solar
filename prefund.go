package solar

import (
	"encoding/json"
	"fmt"
	"strings"
)

func init() {
	cli := app.Command("prefund", "(qtum) fund an owner address with uxtos")

	owner := cli.Arg("owner", "contract name or address to fund").Required().String()

	amount := cli.Arg("amount", "fund an uxto with this amount").Required().Float64()
	multiples := cli.Arg("multiples", "fund this number of identical uxtos").Default("1").Int()

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

		// The JSON object is allowed to have duplicate keys for this call
		// { <addr>: <amount>, ... }

		var uxtos []string
		for i := 0; i < *multiples; i++ {
			uxto := fmt.Sprintf(`"%s": %f`, ownerAddr, *amount)
			uxtos = append(uxtos, uxto)
		}

		amounts := "{\n" + strings.Join(uxtos, ",\n") + "\n}"

		// fmt.Println("jsonuxtos", amounts)

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
