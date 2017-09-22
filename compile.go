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
		var opts CompilerOptions
		return runCompileFile(*file, opts)
	}
}

func runCompileFile(filename string, opts CompilerOptions) error {
	// src, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	return errors.Wrap(err, "Read source")
	// }

	contracts, err := compileSource(filename, opts)
	if err, ok := err.(*CompilerError); ok {
		fmt.Println(err.ErrorOutput)
		os.Exit(1)
	}

	if err != nil {
		return errors.Wrap(err, "compile")
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(contracts)
	if err != nil {
		return err
	}

	return nil
}
