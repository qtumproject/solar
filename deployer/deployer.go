package deployer

import (
	"github.com/qtumproject/solar/contract"
)

type Deployer interface {
	// FIXME better interface for call options
	CreateContract(contract *contract.CompiledContract, jsonParams []byte, name string, overwrite bool, asLib bool, gasLimit int) error
	ConfirmContract(contract *contract.DeployedContract) error
	Mine() error
}
