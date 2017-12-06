package solar

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func init() {
	cmd := app.Command("compile", "Compile Solidity contracts.")
	file := cmd.Arg("file", "Solidity contract source file (.sol)").Required().String()

	appTasks["compile"] = func() (err error) {
		opts, err := solar.SolcOptions()
		if err != nil {
			return
		}

		repo := solar.ContractsRepository()

		c := Compiler{
			Opts:     *opts,
			Filename: *file,
			Repo:     repo,
		}

		contract, err := c.Compile()
		if err, ok := err.(*CompilerError); ok {
			fmt.Println(err.ErrorOutput)
			os.Exit(1)
		}

		if err != nil {
			return errors.Wrap(err, "compile")
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		err = enc.Encode(contract)
		if err != nil {
			return err
		}

		return nil
	}
}
