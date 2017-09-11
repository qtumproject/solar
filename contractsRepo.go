package solar

import (
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
)

type DeployedContracts map[string]*DeployedContract

type DeployedContract struct {
	Name          string `json:"name"`
	Address       Bytes  `json:"address"`
	TransactionID Bytes  `json:"txid"`
	CompiledContract
	CreatedAt time.Time `json:"createdAt"`
	Confirmed bool      `json:"confirmed"`
}

type contractsRepository struct {
	filepath  string
	contracts DeployedContracts
}

func openContractsRepository(filepath string) (repo *contractsRepository, err error) {
	f, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return &contractsRepository{
			filepath:  filepath,
			contracts: make(DeployedContracts),
		}, nil
	}

	if err != nil {
		return
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	contracts := make(DeployedContracts)
	err = dec.Decode(&contracts)
	if err != nil {
		return
	}

	return &contractsRepository{
		filepath:  filepath,
		contracts: contracts,
	}, nil
}

func (r *contractsRepository) Exists(name string) bool {
	_, found := r.contracts[name]
	return found
}

func (r *contractsRepository) Confirm(name string) (err error) {
	c, found := r.contracts[name]
	if !found {
		return errors.Errorf("Cannot unconfirm unknown contract %s", name)
	}

	c.Confirmed = true

	return nil
}

func (r *contractsRepository) Set(name string, c *DeployedContract) (err error) {
	r.contracts[name] = c
	return nil
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

	return enc.Encode(r.contracts)
}
