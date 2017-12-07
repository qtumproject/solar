package solar

import "fmt"

func init() {
	_ = app.Command("confirm", "Wait chain to confirm contracts.")

	appTasks["confirm"] = func() (err error) {
		repo := solar.ContractsRepository()

		err = repo.ConfirmAll(getConfirmUpdateProgressFunc(), solar.Deployer().ConfirmContract)
		if err != nil {
			return
		}

		fmt.Println("All deployed contracts confirmed")
		return
	}
}

func getConfirmUpdateProgressFunc() func (i, total int) {
	reporter := solar.Reporter()
	return func(i , total int) {
		reporter.Submit(eventProgress{
			info: fmt.Sprintf("(%d/%d) Confirming contracts", i, total),
		})

		if i == total {
			reporter.Submit(eventProgressEnd{
				info: fmt.Sprintf("\U0001f680  All contracts confirmed"),
			})
		}
	}
}
