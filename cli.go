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

func init() {
	cmd := app.Command("build", "Compile Solidity contracts.")
	target := cmd.Arg("target", "Source file or directory. Default is the current directory (`.`)").Default(".").String()

	appTasks["build"] = func() (err error) {
		return buildTarget(*target)
	}
}

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
