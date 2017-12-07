package deployer

import (
	"github.com/qtumproject/solar/contract"
)

type Deployer interface {
	CreateContract(contract *contract.CompiledContract, jsonParams []byte, name string, overwrite bool, asLib bool) error
	ConfirmContract(contract *contract.DeployedContract) error
	Mine() error
}
