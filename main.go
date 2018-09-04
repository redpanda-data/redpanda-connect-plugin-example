package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Jeffail/benthos/lib/input"
	"github.com/Jeffail/benthos/lib/log"
	"github.com/Jeffail/benthos/lib/message"
	"github.com/Jeffail/benthos/lib/metrics"
	"github.com/Jeffail/benthos/lib/types"
)

func init() {
	input.RegisterPlugin(
		"example",
		func() interface{} {
			return NewExampleConfig()
		},
		func(iconf interface{}, mgr types.Manager, logger log.Modular, stats metrics.Type) (types.Input, error) {
			return NewExample(iconf, mgr, logger, stats)
		},
	)

	input.DocumentPlugin(
		"example",
		`This plugin example creates an input that generates gibberish messages.`,
		nil, // No need to sanitise the config.
	)
}

// ExampleConfig contains config fields for our plugin type.
type ExampleConfig struct {
	Length int `json:"length" yaml:"length"`
}

// NewExampleConfig creates a config with default values.
func NewExampleConfig() *ExampleConfig {
	return &ExampleConfig{
		Length: 1000,
	}
}

// Example is an example plugin that creates gibberish messages.
type Example struct {
	size int

	transactionsChan chan types.Transaction

	closeChan  chan struct{}
	closedChan chan struct{}
}

// NewExample creates a new example plugin input type.
func NewExample(
	iconf interface{},
	mgr types.Manager,
	log log.Modular,
	stats metrics.Type,
) (input.Type, error) {
	conf, ok := iconf.(*ExampleConfig)
	if !ok {
		return nil, errors.New("failed to cast config")
	}

	e := &Example{
		size: conf.Length,
	}

	go e.loop()
	return e, nil
}

//------------------------------------------------------------------------------

func (e *Example) loop() {
	defer func() {
		close(e.transactionsChan)
		close(e.closedChan)
	}()

	resChan := make(chan types.Response)
	for {
		b := make([]byte, e.size)
		for k := range b {
			b[k] = byte(rand.Int())
		}
		select {
		case e.transactionsChan <- types.NewTransaction(
			message.New([][]byte{b}),
			resChan,
		):
		case <-e.closeChan:
			return
		}
		select {
		case <-resChan:
		case <-e.closeChan:
			return
		}
	}
}

// TransactionChan returns a transactions channel for consuming messages from
// this input type.
func (e *Example) TransactionChan() <-chan types.Transaction {
	return e.transactionsChan
}

// CloseAsync shuts down the input and stops processing requests.
func (e *Example) CloseAsync() {
	close(e.closeChan)
}

// WaitForClose blocks until the input has closed down.
func (e *Example) WaitForClose(timeout time.Duration) error {
	select {
	case <-e.closedChan:
	case <-time.After(timeout):
		return types.ErrTimeout
	}
	return nil
}
