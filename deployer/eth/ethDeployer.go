package eth

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/hayeah/solar/contract"
	"github.com/pkg/errors"
)

type Deployer struct {
	*contract.ContractsRepository
	Account
	client *rpc.Client
}

func NewDeployer(rpcURL *url.URL, repo *contract.ContractsRepository) (*Deployer, error) {
	address, password := "", ""
	if auth := rpcURL.User; auth != nil {
		address = auth.Username()
		password, _ = auth.Password()
		rpcURL.User = nil
	}
	acc := NewAccount(address, password)

	client, err := rpc.Dial(rpcURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "rpc.Dial")
	}

	return &Deployer{
		ContractsRepository: repo,
		client:              client,
		Account:             acc,
	}, nil
}

func (d *Deployer) CreateContract(c *contract.CompiledContract, jsonParams []byte, name string, overwrite bool, aslib bool) (err error) {
	if overwrite {
		if aslib && d.LibExists(name) {
			return errors.Errorf("library name already used: %s", name)
		} else if !aslib && d.Exists(name) {
			return errors.Errorf("contract name already used: %s", name)
		}
	}

	err = d.Unlock(d.client)
	if err != nil {
		return errors.Wrap(err, "eth.Deployer.Unlock")
	}

	// create contract
	// 0x3612c642dfd37cbda70833e44d9d87b65a33cab49468c468be5f0b9b288eba5b
	bin, err := c.ToBytes(jsonParams)
	if err != nil {
		return err
	}

	t := T{
		From: d.Account.Addr,
		Data: "0x" + bin.String(),
		Gas:  4700000,
	}

	fmt.Printf("T: %#v\n", t)
	var txHash string
	err = d.client.Call(&txHash, "eth_sendTransaction", t)
	if err != nil {
		fmt.Println("sendtransaction error", err)
	}
	fmt.Println("txHash:", txHash)

	deployedContract := &contract.DeployedContract{
		CompiledContract: *c,
		Name:             c.Name,
		DeployName:       name,
		TransactionID:    contract.Bytes(txHash),
		CreatedAt:        time.Now(),
	}

	if aslib {
		d.SetLib(name, deployedContract)
	} else {
		d.Set(name, deployedContract)
	}

	err = d.ContractsRepository.Commit()
	if err != nil {
		return
	}

	return nil
}

func (d *Deployer) ConfirmContract(c *contract.DeployedContract) (err error) {
	type txReceipt struct {
		ContractAddress string `json:"contractAddress"`
		TransactionHash string `json:"transactionHash"`
	}

	result := txReceipt{}
	for {
		err = d.client.Call(&result, "eth_getTransactionReceipt", string(c.TransactionID))
		if err != nil {
			fmt.Println("sendtransaction error", err)
		}
		if err != nil {
			return errors.Wrap(err, "eth.Deployer.ConfirmContract")
		}

		if len(result.ContractAddress) != 0 {
			c.Address = contract.Bytes(result.ContractAddress)
			c.Confirmed = true
			fmt.Printf("confirmed result: %#v\n", result)
			break
		}

		nudge := rand.Intn(500)
		time.Sleep(1*time.Second + time.Duration(nudge)*time.Millisecond)
	}

	return nil
}

func (d *Deployer) Mine() (err error) {
	var result interface{}
	err = d.client.Call(&result, "miner_start", 1)
	fmt.Printf("miner_start %#v\n", result)
	return
}

type T struct {
	From     string
	To       string
	Gas      int
	GasPrice *big.Int
	Value    *big.Int
	Data     string
	Nonce    int
}

// MarshalJSON implements the json.Unmarshaler interface.
func (t T) MarshalJSON() ([]byte, error) {
	params := map[string]interface{}{
		"from": t.From,
	}
	if t.To != "" {
		params["to"] = t.To
	}
	if t.Gas > 0 {
		params["gas"] = IntToHex(t.Gas)
	}
	if t.GasPrice != nil {
		params["gasPrice"] = BigToHex(*t.GasPrice)
	}
	if t.Value != nil {
		params["value"] = BigToHex(*t.Value)
	}
	if t.Data != "" {
		params["data"] = t.Data
	}
	if t.Nonce > 0 {
		params["nonce"] = IntToHex(t.Nonce)
	}

	return json.Marshal(params)
}
