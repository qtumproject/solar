package solar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

func buildTarget(target string) (err error) {
	fi, err := os.Stat(target)

	if os.IsNotExist(err) {
		return errors.Errorf("Compile target not found: %s", target)
	}

	if err != nil {
		return
	}

	builder := Builder{
		compilerOpts: CompilerOptions{},
	}

	if fi.IsDir() {
		pat := path.Join(target, "*.sol")
		matches, err := filepath.Glob(pat)
		if err != nil {
			return errors.Wrap(err, "glob")
		}

		for _, filename := range matches {
			fmt.Println("Compiling:", filename)
			err := builder.Compile(filename)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		err := builder.Compile(target)
		if err != nil {
			fmt.Println(err)
		}
	}

	return
}

type Builder struct {
	compilerOpts CompilerOptions
	outputDir    string
}

func (b *Builder) Compile(filename string) error {
	outputFilename := filename + ".json"

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	source, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	compiledContracts, err := compileSource(source, b.compilerOpts)
	if err != nil {
		return errors.Wrap(err, "compile")
	}

	outf, err := os.Create(outputFilename)
	if err != nil {
		return errors.Wrap(err, "output")
	}
	defer outf.Close()

	contracts := make(map[string]CompiledContract)
	for _, contract := range compiledContracts {
		contracts[contract.Name] = contract
	}

	enc := json.NewEncoder(outf)
	enc.SetIndent("", "\t")
	err = enc.Encode(contracts)
	if err != nil {
		return err
	}

	return nil
}
