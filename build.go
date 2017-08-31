package solar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pkg/errors"
)

func init() {
	cmd := app.Command("build", "Compile Solidity contracts.")
	target := cmd.Arg("target", "Source file or directory. Default is the current directory (`.`)").Default(".").String()
	outputDir := cmd.Flag("outdir", "Output directory").String()

	appTasks["build"] = func() (err error) {
		builder := Builder{
			target:       *target,
			outputDir:    *outputDir,
			compilerOpts: CompilerOptions{},
		}

		return builder.build()
	}
}

type Builder struct {
	compilerOpts CompilerOptions
	target       string
	outputDir    string
}

func (b *Builder) build() (err error) {
	target := b.target

	fi, err := os.Stat(target)

	if os.IsNotExist(err) {
		return errors.Errorf("Compile target not found: %s", target)
	}

	if err != nil {
		return
	}

	if fi.IsDir() {
		pat := path.Join(target, "*.sol")
		matches, err := filepath.Glob(pat)
		if err != nil {
			return errors.Wrap(err, "glob")
		}

		limit := make(chan struct{}, runtime.NumCPU())
		var wg sync.WaitGroup
		wg.Add(len(matches))
		for _, filename := range matches {
			limit <- struct{}{}
			filename := filename
			go func() {
				defer func() {
					<-limit
					wg.Done()
				}()
				fmt.Println("Compiling:", filename)
				err := b.compile(filename)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}

		wg.Wait()
	} else {
		err := b.compile(target)
		if err != nil {
			fmt.Println(err)
		}
	}

	return
}

func (b *Builder) compile(filename string) error {
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
