package eth

import (
	"github.com/ethereum/go-ethereum/rpc"
	"fmt"
	"github.com/pkg/errors"
	"bufio"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"os"
)

type Account struct {
	Addr string
	Password string
}

func (acc Account) Unlock(client *rpc.Client) (err error) {
	var result bool

	err = client.Call(&result, "personal_unlockAccount", acc.Addr, acc.Password)
	if err != nil {
		return errors.Wrap(err, "personal_unlockAccount")
	}

	//fmt.Println("unlock account result:", result)
	if result != true {
		return errors.New("unlock account error")
	}

	return nil
}

func NewAccount(addr, password string) Account {
	 acc := Account{addr, password}

	 if len(acc.Addr) == 0 {
		 reader := bufio.NewReader(os.Stdin)

		 fmt.Print("Enter ETH Account: ")
		 username, _ := reader.ReadString('\n')
		 acc.Addr = strings.TrimSpace(username)
	 }
	 if len(acc.Password) == 0 {
		 fmt.Print("Enter Password: ")
		 bytePassword, _ := terminal.ReadPassword(0)
		 acc.Password = strings.TrimSpace(string(bytePassword))
	 }

	 return acc
}