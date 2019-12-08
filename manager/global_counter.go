package manager

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/Jeffail/benthos/v3/lib/log"
	"github.com/Jeffail/benthos/v3/lib/manager"
	"github.com/Jeffail/benthos/v3/lib/metrics"
	"github.com/Jeffail/benthos/v3/lib/types"
	"github.com/Jeffail/gabs/v2"
)

//------------------------------------------------------------------------------

func init() {
	manager.RegisterPlugin(
		"global_counter",
		func() interface{} {
			return NewGlobalCounterConfig()
		},
		func(iconf interface{}, mgr types.Manager, logger log.Modular, stats metrics.Type) (interface{}, error) {
			conf, ok := iconf.(*GlobalCounterConfig)
			if !ok {
				return nil, errors.New("failed to cast config")
			}
			return NewGlobalCounter(*conf)
		},
	)
	manager.DocumentPlugin(
		"global_counter",
		`
Counts messages globally (across all processing threads) and injects the current
count into messages at a given JSON path.`,
		nil, // No need to sanitise the config.
	)
}

// GlobalCounterConfig contains config fields for our plugin type.
type GlobalCounterConfig struct {
	Path string `json:"path" yaml:"path"`
}

// NewGlobalCounterConfig creates a config with default values.
func NewGlobalCounterConfig() *GlobalCounterConfig {
	return &GlobalCounterConfig{
		Path: "counter",
	}
}

//------------------------------------------------------------------------------

// GlobalCounter is an example plugin that creates gibberish messages.
type GlobalCounter struct {
	path    string
	counter uint64
}

// NewGlobalCounter creates a new example plugin input type.
func NewGlobalCounter(conf GlobalCounterConfig) (interface{}, error) {
	return &GlobalCounter{
		path: conf.Path,
	}, nil
}

// Count the parts of a message and injects the running total into each message
// part at a configured path.
func (g *GlobalCounter) Count(msg types.Message) {
	msg.Iter(func(i int, p types.Part) error {
		// Must use atomic as this method can be called from any processing
		// thread.
		total := atomic.AddUint64(&g.counter, 1)
		if jObj, err := p.JSON(); err == nil {
			gObj := gabs.Wrap(jObj)
			gObj.SetP(total, g.path)
			p.SetJSON(gObj.Data())
		}
		return nil
	})
}

//------------------------------------------------------------------------------

// CloseAsync shuts down the plugin.
func (g *GlobalCounter) CloseAsync() {
	// No need to clean up any resources.
}

// WaitForClose blocks until the plugin has closed down.
func (g *GlobalCounter) WaitForClose(timeout time.Duration) error {
	// Nothing to do here.
	return nil
}

//------------------------------------------------------------------------------
