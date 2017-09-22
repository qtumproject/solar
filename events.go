package solar

import (
	"fmt"
	"time"

	spin "github.com/tj/go-spin"
)

type eventReporter func(interface{})

type events struct {
	in        chan interface{}
	reporters []eventReporter
}

func (r *events) Submit(item interface{}) {
	r.in <- item
}

// func (r *events) Subscribe(fn eventReporter) {
// 	r.reporters = append(r.reporters, fn)
// }

func (r *events) Start() {
	ticker := time.NewTicker(150 * time.Millisecond)
	var tickerInfo string

	spinner := spin.New()
	spinner.Set(spin.Spin3)

	for {
		select {
		case <-ticker.C:
			if tickerInfo != "" {
				fmt.Printf("\r%s %s", spinner.Next(), tickerInfo)
			}
		case e := <-r.in:
			switch e := e.(type) {
			case eventProgress:
				tickerInfo = e.info
				// log.Println(e)
			case eventProgressEnd:
				fmt.Printf("\033[2K\033[1000D%s\n", e.info)
				tickerInfo = ""
				// log.Println(e)
			}
		}

		// for _, reporterFn := range r.reporters {
		// 	reporterFn(event)
		// }

	}
}

type eventProgress struct {
	// [spinner] [info]
	info string
}

type eventProgressEnd struct {
	info string
}

// type eventConfirmContract struct {
// }

// type eventConfirmContractsStart struct {
// 	total int
// }

// type contractConfirmationReporter struct {
// 	init sync.Once
// }

// func (r *contractConfirmationReporter) Report(e interface{}) (err error) {
// 	switch e := e.(type) {
// 	case *eventConfirmContractsStart:

// 	}

// 	// if e, ok := i.(*eventConfirmContractProgress); ok {
// 	// 	// if e.n == 0 {
// 	// 	// 	r.init.Do(func() {

// 	// 	// 	})
// 	// 	// }
// 	// 	fmt.Sprintf("[%d/%d] Confirming contracts")
// 	// }
// }

// func (r *contractConfirmationReporter) Start() {
// }
