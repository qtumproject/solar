package solar

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// {"version": "1.1", "method": "confirmFruitPurchase", "params": [["apple", "orange", "mangoes"], 1.123], "id": "194521489"}

type JSONRPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     interface{}   `json:"id"`
}

func uploadContract(contract *CompiledContract, gasLimit int) (err error) {
	// qtumd
	url := "http://localhost:3889/"

	user := "howard"
	password := "yeh"

	// jsonReq := JSONRPCRequest{
	// 	Method: "getaccountinfo",
	// 	Params: []interface{}{
	// 		// "142eea127133fb5c9f2d10d10559753d9a968475",
	// 		"142eea127133fb5c9f2d10d10559753d9a968475",
	// 	},
	// }

	jsonReq := JSONRPCRequest{
		Method: "createcontract",
		Params: []interface{}{
			contract.Bin.String(),
			gasLimit,
		},
	}

	var body bytes.Buffer
	enc := json.NewEncoder(&body)
	err = enc.Encode(&jsonReq)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return
	}
	userpass := fmt.Sprintf("%s:%s", user, password)
	auth := base64.RawStdEncoding.EncodeToString([]byte(userpass))

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	log.Println("status", res.Status)

	_, err = io.Copy(os.Stderr, res.Body)

	return
}
