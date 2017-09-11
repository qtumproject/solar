package solar

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func init() {
	_ = app.Command("confirm", "Wait for contract creation to complete.")

	appTasks["confirm"] = func() (err error) {
		repo := solar.ContractsRepository()
		rpc := solar.RPC()

		var wg sync.WaitGroup
		wg.Add(len(repo.contracts))
		for name, contract := range repo.contracts {
			if contract.Confirmed {
				wg.Done()
				continue
			}

			name := name
			contract := contract
			go func() {
				err := confirmDeployedContract(rpc, name, contract)
				if err != nil {
					log.Println("err", err)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		err = repo.Commit()
		if err != nil {
			return
		}

		fmt.Println("All contracts deployment confirmed")

		return
	}
}

func confirmDeployedContract(rpc *qtumRPC, name string, c *DeployedContract) (err error) {
	for {
		fmt.Printf("Checking %s\n", name)
		result := make(map[string]interface{})
		err := rpc.Call(&result, "getaccountinfo", c.Address)
		if err, ok := err.(*jsonRPCError); ok {
			fmt.Printf("%s\t%s\n", name, err)
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
