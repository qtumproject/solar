package solar

func init() {
	_ = app.Command("confirm", "Wait chain to confirm contracts.")

	appTasks["confirm"] = func() (err error) {
		repo := solar.ContractsRepository()
		err = repo.ConfirmAll()
		if err != nil {
			return
		}

		// fmt.Println("All deployed contracts confirmed")
		return
	}
}
