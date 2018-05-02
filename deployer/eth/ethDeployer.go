package eth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net/url"
	"time"

	"github.com/qtumproject/solar/deployer"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/qtumproject/solar/contract"
)

type Deployer struct {
	*contract.ContractsRepository
	Account
	client *rpc.Client
}

func NewDeployer(rpcURL *url.URL, repo *contract.ContractsRepository) (*Deployer, error) {
	auth := rpcURL.User
	if auth == nil {
		return nil, errors.New("address and password not specified")
	}

	address := auth.Username()
	password, _ := auth.Password()

	// acc := NewAccount(address, password)

	client, err := rpc.Dial(rpcURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "rpc.Dial")
	}

	return &Deployer{
		ContractsRepository: repo,
		client:              client,
		Account: Account{
			Addr:     address,
			Password: password,
		},
	}, nil
}

func (d *Deployer) CreateContract(c *contract.CompiledContract, jsonParams []byte, opts *deployer.Options) (err error) {
	name := opts.Name

	if !opts.Overwrite {
		if opts.AsLib && d.LibExists(name) {
			return errors.Errorf("library name already used: %s", name)
		} else if !opts.AsLib && d.Exists(name) {
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

	var gasPrice big.Float // gwei
	gasPrice.Copy(opts.GasPrice)
	gasPrice.Mul(&gasPrice, big.NewFloat(1e9))

	var gasPriceWei big.Int
	gasPrice.Int(&gasPriceWei)

	fmt.Println("gasPrice", opts.GasPrice.String(), gasPrice.String())
	fmt.Println("gasPriceWei", gasPriceWei.String())

	t := T{
		From:     d.Account.Addr,
		Data:     bin.String(),
		Gas:      big.NewInt(int64(opts.GasLimit)),
		GasPrice: &gasPriceWei,
	}

	//fmt.Printf("T: %#v\n", t)
	var txHash string
	err = d.client.Call(&txHash, "eth_sendTransaction", t)
	if err != nil {
		fmt.Println("sendtransaction error", err)
		return errors.Wrap(err, "sendtransaction")
	}
	fmt.Printf("txHash: %s\n", txHash)

	hexBytes, _ := hex.DecodeString(txHash[2:])
	deployedContract := &contract.DeployedContract{
		CompiledContract: *c,
		Name:             c.Name,
		DeployName:       name,
		TransactionID:    contract.Bytes(hexBytes),
		CreatedAt:        time.Now(),
		Sender:           d.Account.Addr,
	}

	if opts.AsLib {
		d.ContractsRepository.SetLib(name, deployedContract)
	} else {
		d.ContractsRepository.Set(name, deployedContract)
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
		err = d.client.Call(&result, "eth_getTransactionReceipt", "0x"+hex.EncodeToString(c.TransactionID))
		if err != nil {
			fmt.Println("sendtransaction error", err)
		}
		if err != nil {
			return errors.Wrap(err, "eth.Deployer.ConfirmContract")
		}

		if len(result.ContractAddress) != 0 {
			addressBytes, _ := hex.DecodeString(result.ContractAddress[2:])
			c.Address = contract.Bytes(addressBytes)
			c.Confirmed = true
			fmt.Printf("\rcontractAddress: %s\n", result.ContractAddress)
			break
		}

		nudge := rand.Intn(500)
		time.Sleep(1*time.Second + time.Duration(nudge)*time.Millisecond)
	}

	return nil
}

func (d *Deployer) Mine() (err error) {
	// var result interface{}
	// err = d.client.Call(&result, "miner_start", 1)
	//fmt.Printf("miner_start %#v\n", result)
	return
}

type T struct {
	From     string
	To       string
	Gas      *big.Int
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
	if t.Gas.Sign() > 0 {
		params["gas"] = BigToHex(t.Gas)
	}
	if t.GasPrice.Sign() > 0 {
		params["gasPrice"] = BigToHex(t.GasPrice)
	}
	if t.Value != nil {
		params["value"] = BigToHex(t.Value)
	}
	if t.Data != "" {
		params["data"] = t.Data
	}
	if t.Nonce > 0 {
		params["nonce"] = IntToHex(t.Nonce)
	}

	return json.Marshal(params)
}
