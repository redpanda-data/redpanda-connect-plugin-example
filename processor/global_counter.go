package processor

import (
	"errors"
	"fmt"
	"time"

	"github.com/Jeffail/benthos/v3/lib/log"
	"github.com/Jeffail/benthos/v3/lib/metrics"
	"github.com/Jeffail/benthos/v3/lib/processor"
	"github.com/Jeffail/benthos/v3/lib/types"
	"github.com/benthosdev/benthos-plugin-example/manager"
)

//------------------------------------------------------------------------------

func init() {
	processor.RegisterPlugin(
		"global_counter",
		func() interface{} {
			return NewGlobalCounterConfig()
		},
		func(
			iconf interface{},
			mgr types.Manager,
			logger log.Modular,
			stats metrics.Type,
		) (types.Processor, error) {
			conf, ok := iconf.(*GlobalCounterConfig)
			if !ok {
				return nil, errors.New("failed to cast config")
			}
			return NewGlobalCounter(*conf, mgr)
		},
	)
	processor.DocumentPlugin(
		"global_counter",
		`
Applies a global counter to messages. Requires a `+"`global_counter`"+` resource
which can be referenced by any number of `+"`global_counter`"+` processors.`,
		nil,
	)
}

// GlobalCounterConfig contains config fields for our plugin type.
type GlobalCounterConfig struct {
	Resource string `json:"resource" yaml:"resource"`
}

// NewGlobalCounterConfig creates a config with default values.
func NewGlobalCounterConfig() *GlobalCounterConfig {
	return &GlobalCounterConfig{
		Resource: "",
	}
}

//------------------------------------------------------------------------------

// GlobalCounter is a processor that reverses all messages.
type GlobalCounter struct {
	resource *manager.GlobalCounter
}

// NewGlobalCounter returns a GlobalCounter processor.
func NewGlobalCounter(conf GlobalCounterConfig, mgr types.Manager) (types.Processor, error) {
	resource, err := mgr.GetPlugin(conf.Resource)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain resource '%v': %v", conf.Resource, err)
	}

	ctr, ok := resource.(*manager.GlobalCounter)
	if !ok {
		return nil, fmt.Errorf("resource '%v' was an unexpected type", conf.Resource)
	}

	return &GlobalCounter{
		resource: ctr,
	}, nil
}

// ProcessMessage applies the processor to a message
func (m *GlobalCounter) ProcessMessage(msg types.Message) ([]types.Message, types.Response) {
	// Always create a new copy if we intend to mutate message contents.
	newMsg := msg.Copy()
	m.resource.Count(newMsg)
	return []types.Message{newMsg}, nil
}

// CloseAsync shuts down the processor and stops processing requests.
func (m *GlobalCounter) CloseAsync() {
}

// WaitForClose blocks until the processor has closed down.
func (m *GlobalCounter) WaitForClose(timeout time.Duration) error {
	return nil
}

//------------------------------------------------------------------------------
