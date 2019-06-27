package qtum

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/qtumproject/solar/contract"
)

type jsonRPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     string        `json:"id"`
}

type jsonRPCRersult struct {
	RawResult json.RawMessage `json:"result"`
	// RawError  json.RawMessage `json:"error"`
	Error *jsonRPCError `json:"error"`
	ID    string        `json:"id"`
}

type jsonRPCError struct {
	Code    int
	Message string
}

func (err *jsonRPCError) Error() string {
	return fmt.Sprintf("[code: %d] %s", err.Code, err.Message)
}

type TransactionReceipt struct {
	TxID    contract.Bytes `json:"txid"`
	Sender  string         `json:"sender"`
	Hash160 contract.Bytes `json:"hash160"`
	Address contract.Bytes `json:"address"`
}

func NewRPC(baseurl string) (*RPC, error) {
	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, err
	}

	return &RPC{
		BaseURL: u,
	}, nil
}

type RPC struct {
	BaseURL *url.URL
}

func (rpc *RPC) Call(result interface{}, method string, params ...interface{}) (err error) {
	url := rpc.BaseURL

	jsonReq := jsonRPCRequest{
		Method: method,
		Params: params,
	}

	var body bytes.Buffer
	enc := json.NewEncoder(&body)
	err = enc.Encode(&jsonReq)
	if err != nil {
		return
	}

	// would user info be included in the URL?
	urlString := url.String()
	// log.Println("rpc url", urlString)

	req, err := http.NewRequest(http.MethodPost, urlString, &body)
	if err != nil {
		return
	}

	if auth := url.User; auth != nil {
		password, _ := auth.Password()
		userpass := fmt.Sprintf("%s:%s", auth.Username(), password)
		token := base64.RawStdEncoding.EncodeToString([]byte(userpass))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	// log.Println("rpc http status", res.Status)
	if result == nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("QTUM RPC %s", res.Status)
	}

	dec := json.NewDecoder(res.Body)
	jsonResult := &jsonRPCRersult{}
	err = dec.Decode(jsonResult)
	if err != nil {
		return
	}

	if res.StatusCode == 200 {
		// pretty.Println("json rpc result:", string(jsonResult.RawResult))
		json.Unmarshal(jsonResult.RawResult, result)
		return
	}

	// QTum RPC returns 500 for RPC error
	return jsonResult.Error
}
