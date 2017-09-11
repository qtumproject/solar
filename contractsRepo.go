package solar

import (
	"encoding/json"
	"os"
	"time"
)

type DeployedContracts map[string]DeployedContract

type DeployedContract struct {
	Name          string `json:"name"`
	Address       Bytes  `json:"address"`
	TransactionID Bytes  `json:"txid"`
	CompiledContract
	CreatedAt time.Time `json:"createdAt"`
}

type deployedContractsRepository struct {
	filepath  string
	contracts DeployedContracts
}

func openDeployedContractsRepository(filepath string) (repo *deployedContractsRepository, err error) {
	f, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return &deployedContractsRepository{
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

	return &deployedContractsRepository{
		filepath:  filepath,
		contracts: contracts,
	}, nil
}

func (r *deployedContractsRepository) Add(name string, c DeployedContract) (err error) {
	// TODO check if contract already exists
	r.contracts[name] = c
	return nil
}

func (r *deployedContractsRepository) Commit() (err error) {
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
