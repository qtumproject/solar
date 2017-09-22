package solar

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app       = kingpin.New("solar", "Solidity smart contract deployment management.")
	solarEnv  = app.Flag("env", "Environment name").Envar("SOLAR_ENV").Default("development").String()
	solarRPC  = app.Flag("rpc", "RPC provider url").Envar("SOLAR_RPC").String()
	solarRepo = app.Flag("repo", "Path of contracts repository").Envar("SOLAR_REPO").String()
	appTasks  = map[string]func() error{}
)

type solarCLI struct {
	rpc     *qtumRPC
	rpcOnce sync.Once

	repo     *contractsRepository
	repoOnce sync.Once

	reporter     *events
	reporterOnce sync.Once
}

var solar = &solarCLI{
// reporter: &events{
// 	in: make(chan interface{}),
// }
}

func (c *solarCLI) Reporter() *events {
	c.reporterOnce.Do(func() {
		c.reporter = &events{
			in: make(chan interface{}),
		}

		go c.reporter.Start()
	})

	return c.reporter
}

func (c *solarCLI) RPC() *qtumRPC {
	c.rpcOnce.Do(func() {
		rpcURL, err := url.Parse(*solarRPC)
		if err != nil {
			fmt.Println("Invalid RPC URL:", rpcURL)
			os.Exit(1)
		}

		c.rpc = &qtumRPC{rpcURL}
	})

	return c.rpc
}

// Open the file `solar.{SOLAR_ENV}.json` as contracts repository
func (c *solarCLI) ContractsRepository() *contractsRepository {
	c.repoOnce.Do(func() {
		var repoFilePath string
		if *solarRepo != "" {
			repoFilePath = *solarRepo
		} else {
			repoFilePath = fmt.Sprintf("solar.%s.json", *solarEnv)
		}

		repo, err := openContractsRepository(repoFilePath)
		if err != nil {
			fmt.Println("Cannot open contracts repo:", repoFilePath)
			os.Exit(1)
		}

		c.repo = repo
	})

	return c.repo
}

func (c *solarCLI) Deployer() *Deployer {
	return &Deployer{
		rpc:  c.RPC(),
		repo: c.ContractsRepository(),
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
		fmt.Println(err)
		os.Exit(1)
	}
}