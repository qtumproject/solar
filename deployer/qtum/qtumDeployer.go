package qtum

import (
	"net/url"
	"time"

	"github.com/qtumproject/solar/deployer"

	"github.com/qtumproject/solar/b58addr"

	"math/rand"

	"github.com/pkg/errors"
	"github.com/qtumproject/solar/contract"
)

type Deployer struct {
	rpc *RPC
	*contract.ContractsRepository

	// qtum base58 sender address used to create a contract.
	senderAddress string
}

func NewDeployer(rpcURL *url.URL, repo *contract.ContractsRepository, senderAddress string) (*Deployer, error) {
	return &Deployer{
		rpc: &RPC{
			BaseURL: rpcURL,
		},
		ContractsRepository: repo,
		senderAddress:       senderAddress,
	}, nil
}

func (d *Deployer) Mine() error {
	return d.rpc.Call(nil, "generate", 1)
}

func (d *Deployer) ConfirmContract(c *contract.DeployedContract) (err error) {
	for {
		// fmt.Printf("Checking %s\n", name)
		result := make(map[string]interface{})
		err := d.rpc.Call(&result, "getaccountinfo", c.Address)
		if err, ok := err.(*jsonRPCError); ok {
			// fmt.Printf("%s\t%s\n", name, err)
			nudge := rand.Intn(500)
			time.Sleep(1*time.Second + time.Duration(nudge)*time.Millisecond)
			continue
		} else if err != nil {
			return err
		}

		// fmt.Printf("confirmed\t%s\t%s\n", name, c.Address)
		c.Confirmed = true
		return nil
	}
}

func (d *Deployer) CreateContract(c *contract.CompiledContract, jsonParams []byte, opts *deployer.Options) (err error) {
	// TODO: dry out similar CreateContract code from eth and qtum deployers
	name := opts.Name

	if !opts.Overwrite {
		if opts.AsLib && d.LibExists(name) {
			return errors.Errorf("library name already used: %s", name)
		} else if !opts.AsLib && d.Exists(name) {
			return errors.Errorf("contract name already used: %s", name)
		}
	}

	bin, err := c.ToBytes(jsonParams)
	if err != nil {
		return
	}

	var tx TransactionReceipt

	var gasLimit uint

	if opts.GasLimit > 0 {
		gasLimit = opts.GasLimit
	} else {
		gasLimit = 300000
	}

	args := []interface{}{
		bin, gasLimit, 0.0000004,
	}

	// fmt.Println("create contract args", args)

	if d.senderAddress != "" {
		args = append(args, d.senderAddress)
	}

	err = d.rpc.Call(&tx, "createcontract", args...)

	if err != nil {
		return errors.Wrap(err, "createcontract")
	}

	// fmt.Println("tx", tx.Address)
	// fmt.Println("contract name", contract.Name)

	deployedContract := &contract.DeployedContract{
		CompiledContract: *c,
		Name:             c.Name,
		DeployName:       name,
		TransactionID:    tx.TxID,
		Address:          tx.Address,
		CreatedAt:        time.Now(),
		Sender:           tx.Sender,
		SenderHex:        b58addr.ToHexString(tx.Sender),
	}

	if opts.AsLib {
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
