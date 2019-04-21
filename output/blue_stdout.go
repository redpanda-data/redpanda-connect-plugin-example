package output

import (
	"fmt"
	"sync"
	"time"

	"github.com/Jeffail/benthos/lib/log"
	"github.com/Jeffail/benthos/lib/metrics"
	"github.com/Jeffail/benthos/lib/output"
	"github.com/Jeffail/benthos/lib/response"
	"github.com/Jeffail/benthos/lib/types"
)

func init() {
	output.RegisterPlugin(
		"blue_stdout",
		func() interface{} {
			conf := NewBlueStdoutConfig()
			return &conf
		},
		func(iconf interface{}, mgr types.Manager, logger log.Modular, stats metrics.Type) (types.Output, error) {
			return NewBlueStdout(mgr, logger, stats)
		},
	)

	output.DocumentPlugin(
		"blue_stdout",
		`
This plugin example creates an output that writes messages to stdout BUT THE
TEXT IS HECKING BLUE!`,
		nil, // No need to sanitise the config.
	)
}

//------------------------------------------------------------------------------

// BlueStdoutConfig contains configuration fields for the BlueStdout output.
type BlueStdoutConfig struct {
}

// NewBlueStdoutConfig returns a BlueStdoutConfig with default values.
func NewBlueStdoutConfig() BlueStdoutConfig {
	return BlueStdoutConfig{}
}

//------------------------------------------------------------------------------

// BlueStdout is an example plugin that creates gibberish messages.
type BlueStdout struct {
	transactionsChan <-chan types.Transaction

	log   log.Modular
	stats metrics.Type

	closeOnce  sync.Once
	closeChan  chan struct{}
	closedChan chan struct{}
}

// NewBlueStdout creates a new example plugin output type.
func NewBlueStdout(
	mgr types.Manager,
	log log.Modular,
	stats metrics.Type,
) (output.Type, error) {
	e := &BlueStdout{
		log:   log,
		stats: stats,

		closeChan:  make(chan struct{}),
		closedChan: make(chan struct{}),
	}

	return e, nil
}

//------------------------------------------------------------------------------

func (e *BlueStdout) loop() {
	defer func() {
		close(e.closedChan)
	}()

	for {
		var tran types.Transaction
		var open bool
		select {
		case tran, open = <-e.transactionsChan:
			if !open {
				return
			}
		case <-e.closeChan:
			return
		}
		tran.Payload.Iter(func(i int, p types.Part) error {
			fmt.Printf("\033[01;34m%s\033[m\n", p.Get())
			return nil
		})
		select {
		case tran.ResponseChan <- response.NewAck():
		case <-e.closeChan:
			return
		}
	}
}

// Connected returns true if this output is currently connected to its target.
func (e *BlueStdout) Connected() bool {
	return true // We're always connected
}

// Consume starts this output consuming from a transaction channel.
func (e *BlueStdout) Consume(tChan <-chan types.Transaction) error {
	e.transactionsChan = tChan
	go e.loop()
	return nil
}

// CloseAsync shuts down the output and stops processing requests.
func (e *BlueStdout) CloseAsync() {
	e.closeOnce.Do(func() {
		close(e.closeChan)
	})
}

// WaitForClose blocks until the output has closed down.
func (e *BlueStdout) WaitForClose(timeout time.Duration) error {
	select {
	case <-e.closedChan:
	case <-time.After(timeout):
		return types.ErrTimeout
	}
	return nil
}
