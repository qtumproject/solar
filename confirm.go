package solar

import (
	"fmt"
)

func init() {
	_ = app.Command("confirm", "Wait for contract creation to complete.")

	appTasks["confirm"] = func() (err error) {
		repo := solar.ContractsRepository()
		err = repo.ConfirmAll()
		if err != nil {
			return
		}

		fmt.Println("All deployed contracts confirmed")
		return
	}
}
