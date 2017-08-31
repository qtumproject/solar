package solar

import (
	"log"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("solar", "Solidity smart contract deployment management.")
	appTasks = map[string]func() error{}
)

func Main() {
	cmdName, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	task := appTasks[cmdName]
	err = task()
	if err != nil {
		log.Fatalln(err)
	}
}
