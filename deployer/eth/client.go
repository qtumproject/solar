package eth

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type Account struct {
	Addr     string
	Password string
}

func (acc Account) Unlock(client *rpc.Client) (err error) {
	var result bool

	err = client.Call(&result, "personal_unlockAccount", acc.Addr, acc.Password, nil)
	if err != nil {
		return errors.Wrap(err, "personal_unlockAccount")
	}

	//fmt.Println("unlock account result:", result)
	if result != true {
		return errors.New("unlock account error")
	}

	return nil
}
