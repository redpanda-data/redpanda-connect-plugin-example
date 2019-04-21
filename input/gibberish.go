package input

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/Jeffail/benthos/lib/input"
	"github.com/Jeffail/benthos/lib/log"
	"github.com/Jeffail/benthos/lib/message"
	"github.com/Jeffail/benthos/lib/metrics"
	"github.com/Jeffail/benthos/lib/types"
)

func init() {
	input.RegisterPlugin(
		"gibberish",
		func() interface{} {
			return NewGibberishConfig()
		},
		func(iconf interface{}, mgr types.Manager, logger log.Modular, stats metrics.Type) (types.Input, error) {
			conf, ok := iconf.(*GibberishConfig)
			if !ok {
				return nil, errors.New("failed to cast config")
			}
			return NewGibberish(*conf, mgr, logger, stats)
		},
	)

	input.DocumentPlugin(
		"gibberish",
		`
This plugin example creates an input that generates gibberish ascii messages.`,
		nil, // No need to sanitise the config.
	)
}

// GibberishConfig contains config fields for our plugin type.
type GibberishConfig struct {
	Length int `json:"length" yaml:"length"`
}

// NewGibberishConfig creates a config with default values.
func NewGibberishConfig() *GibberishConfig {
	return &GibberishConfig{
		Length: 1000,
	}
}

// Gibberish is an example plugin that creates gibberish messages.
type Gibberish struct {
	size int

	transactionsChan chan types.Transaction

	log   log.Modular
	stats metrics.Type

	closeOnce  sync.Once
	closeChan  chan struct{}
	closedChan chan struct{}
}

// NewGibberish creates a new example plugin input type.
func NewGibberish(
	conf GibberishConfig,
	mgr types.Manager,
	log log.Modular,
	stats metrics.Type,
) (input.Type, error) {
	e := &Gibberish{
		size: conf.Length,

		log:   log,
		stats: stats,

		transactionsChan: make(chan types.Transaction),
		closeChan:        make(chan struct{}),
		closedChan:       make(chan struct{}),
	}

	go e.loop()
	return e, nil
}

//------------------------------------------------------------------------------

func (e *Gibberish) loop() {
	defer func() {
		close(e.transactionsChan)
		close(e.closedChan)
	}()

	resChan := make(chan types.Response)
	for {
		b := make([]byte, e.size)
		for k := range b {
			b[k] = byte((rand.Int() % 94) + 32)
		}

		// send batch to downstream processors
		select {
		case e.transactionsChan <- types.NewTransaction(
			message.New([][]byte{b}),
			resChan,
		):
		case <-e.closeChan:
			return
		}

		// check transaction success
		select {
		case result := <-resChan:
			if nil != result.Error() {
				e.log.Errorln(result.Error().Error())
				continue
			}
		case <-e.closeChan:
			return
		}
	}
}

// Connected returns true if this input is currently connected to its target.
func (e *Gibberish) Connected() bool {
	return true // We're always connected
}

// TransactionChan returns a transactions channel for consuming messages from
// this input type.
func (e *Gibberish) TransactionChan() <-chan types.Transaction {
	return e.transactionsChan
}

// CloseAsync shuts down the input and stops processing requests.
func (e *Gibberish) CloseAsync() {
	e.closeOnce.Do(func() {
		close(e.closeChan)
	})
}

// WaitForClose blocks until the input has closed down.
func (e *Gibberish) WaitForClose(timeout time.Duration) error {
	select {
	case <-e.closedChan:
	case <-time.After(timeout):
		return types.ErrTimeout
	}
	return nil
}
