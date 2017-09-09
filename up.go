package solar

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kr/pretty"
)

func upTask(compiledSolFiles []string) {

}

// {"version": "1.1", "method": "confirmFruitPurchase", "params": [["apple", "orange", "mangoes"], 1.123], "id": "194521489"}

type JSONRPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     string        `json:"id"`
}

type JSONRPCRersult struct {
	RawResult json.RawMessage `json:"result"`
	RawError  json.RawMessage `json:"error"`
	ID        string          `json:"id"`
}

type TransactionReceipt struct {
	TxID    Bytes  `json:"txid"`
	Sender  string `json:"sender"`
	Hash160 Bytes  `json:"hash160"`
	Address Bytes  `json:"address"`
}

type uploader struct {
	opts      uploaderOptions
	contracts []*CompiledContract
}

type uploaderOptions struct {
	GasLimit uint
	GasPrice uint
}

func uploadContract(contract *CompiledContract, gasLimit, gasPrice int) (err error) {
	// qtumd

	// jsonReq := JSONRPCRequest{
	// 	Method: "getaccountinfo",
	// 	Params: []interface{}{
	// 		// "142eea127133fb5c9f2d10d10559753d9a968475",
	// 		"142eea127133fb5c9f2d10d10559753d9a968475",
	// 	},
	// }

	// jsonReq := JSONRPCRequest{
	// 	Method: "createcontract",
	// 	Params: []interface{}{
	// 		contract.Bin.String(),
	// 		gasLimit,
	// 	},
	// }
	res, err := callRPC("createcontract", contract.Bin.String(), gasLimit)
	if err != nil {
		return
	}

	var tx TransactionReceipt
	json.Unmarshal(res.RawResult, &tx)

	// _, err = io.Copy(os.Stderr, res.Body)
	log.Println("tx", tx)

	// For ragtest, generate one.
	_, err = callRPC("generate", 1)
	if err != nil {
		log.Println("generate 1 block", err)
	}

	for {
		log.Println("look up account")
		res, err := callRPC("getaccountinfo", tx.Address.String())
		if err != nil {
			return err
		}

		log.Println("getaccountinfo", string(res.RawResult))
		if string(res.RawResult) != "null" {
			break
		}

		time.Sleep(1 * time.Second)
	}

	// loop keep looping to look up transaction

	return
}

func callRPC(method string, params ...interface{}) (jsonResult *JSONRPCRersult, err error) {
	url := "http://localhost:13889/"
	user := "howard"
	password := "yeh"

	jsonReq := JSONRPCRequest{
		Method: method,
		Params: params,
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

	log.Println("rpc http status", res.Status)

	dec := json.NewDecoder(res.Body)
	jsonResult = &JSONRPCRersult{}
	err = dec.Decode(jsonResult)
	if err != nil {
		return
	}

	if res.StatusCode == 200 {
		pretty.Println("json rpc result:", string(jsonResult.RawResult))
	} else {
		pretty.Println("json rpc result:", string(jsonResult.RawError))
	}

	return
}
