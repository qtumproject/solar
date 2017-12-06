package solar

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type DeployedContracts map[string]*DeployedContract

type DeployedContract struct {
	Name          string `json:"name"`
	DeployName    string `json:"deployName"`
	Address       Bytes  `json:"address"`
	TransactionID Bytes  `json:"txid"`
	CompiledContract
	CreatedAt time.Time `json:"createdAt"`
	Confirmed bool      `json:"confirmed"`
}

type contractsRepository struct {
	filepath  string
	Contracts DeployedContracts `json:"contracts"`
	Libraries DeployedContracts `json:"libraries"`
}

func openContractsRepository(filepath string) (repo *contractsRepository, err error) {
	f, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return &contractsRepository{
			filepath:  filepath,
			Contracts: make(DeployedContracts),
			Libraries: make(DeployedContracts),
		}, nil
	}

	if err != nil {
		return
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	repo = &contractsRepository{
		filepath: filepath,
	}

	err = dec.Decode(&repo)
	if err != nil {
		return
	}

	if repo.Libraries == nil {
		repo.Libraries = make(DeployedContracts)
	}

	if repo.Contracts == nil {
		repo.Contracts = make(DeployedContracts)
	}

	return
}

func (r *contractsRepository) UnconfirmedContracts() []*DeployedContract {
	var contracts []*DeployedContract

	for _, contract := range r.Contracts {
		if !contract.Confirmed {
			contracts = append(contracts, contract)
		}
	}

	return contracts
}

func (r *contractsRepository) SortedContracts() []*DeployedContract {
	var contracts []*DeployedContract

	for _, contract := range r.Contracts {
		contracts = append(contracts, contract)
	}

	sort.Slice(contracts, func(i, j int) bool {
		c1 := contracts[i]
		c2 := contracts[j]

		return c1.CreatedAt.Unix() < c2.CreatedAt.Unix()
	})

	return contracts
}

func (r *contractsRepository) Exists(name string) bool {
	_, found := r.Contracts[name]
	return found
}

func (r *contractsRepository) LibExists(name string) bool {
	_, found := r.Libraries[name]
	return found
}

func (r *contractsRepository) Confirm(name string) (err error) {
	c, found := r.Contracts[name]
	if !found {
		return errors.Errorf("Cannot unconfirm unknown contract %s", name)
	}

	c.Confirmed = true

	return nil
}

func (r *contractsRepository) Set(name string, c *DeployedContract) {
	r.Contracts[name] = c
}

func (r *contractsRepository) SetLib(name string, c *DeployedContract) {
	r.Libraries[name] = c
}

func (r *contractsRepository) Commit() (err error) {
	// TODO. do write & swap instead of truncat?
	f, err := os.Create(r.filepath)
	if err != nil {
		return
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	return enc.Encode(r)
}

// Confirm checks the RPC server to see if all the contracts
// are confirmed by the blockchain.
func (r *contractsRepository) ConfirmAll() (err error) {
	contracts := r.UnconfirmedContracts()

	total := len(contracts)

	reporter := solar.Reporter()

	updateProgress := func(i int) {
		reporter.Submit(eventProgress{
			info: fmt.Sprintf("(%d/%d) Confirming contracts", i, total),
		})

		if i == total {
			reporter.Submit(eventProgressEnd{
				info: fmt.Sprintf("\U0001f680  All contracts confirmed"),
			})
		}

	}

	updateProgress(0)

	for i, contract := range contracts {
		contract := contract

		err := r.confirmContract(contract)
		if err != nil {
			log.Println("err", err)
		}

		updateProgress(i + 1)
	}

	err = r.Commit()
	if err != nil {
		return
	}

	return
}

func (r *contractsRepository) confirmContract(c *DeployedContract) (err error) {
	rpc := solar.RPC()

	// name := c.Name
	for {
		// fmt.Printf("Checking %s\n", name)
		result := make(map[string]interface{})
		err := rpc.Call(&result, "getaccountinfo", c.Address)
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
