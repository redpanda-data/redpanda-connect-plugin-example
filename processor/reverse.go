package processor

import (
	"time"

	"github.com/Jeffail/benthos/lib/log"
	"github.com/Jeffail/benthos/lib/metrics"
	"github.com/Jeffail/benthos/lib/processor"
	"github.com/Jeffail/benthos/lib/types"
)

//------------------------------------------------------------------------------

func init() {
	processor.RegisterPlugin(
		"reverse", nil, // No configuration needed
		func(
			iconf interface{},
			mgr types.Manager,
			logger log.Modular,
			stats metrics.Type,
		) (types.Processor, error) {
			return NewReverse(logger, stats)
		},
	)
	processor.DocumentPlugin(
		"reverse",
		`Reverses the raw bytes of every message.`,
		nil,
	)
}

//------------------------------------------------------------------------------

// Reverse is a processor that reverses all messages.
type Reverse struct {
	log   log.Modular
	stats metrics.Type
}

// NewReverse returns a Reverse processor.
func NewReverse(
	log log.Modular, stats metrics.Type,
) (types.Processor, error) {
	m := &Reverse{
		log:   log,
		stats: stats,
	}
	return m, nil
}

// ProcessMessage applies the processor to a message
func (m *Reverse) ProcessMessage(msg types.Message) ([]types.Message, types.Response) {
	// Always create a new copy if we intend to mutate message contents.
	newMsg := msg.Copy()
	newMsg.Iter(func(i int, p types.Part) error {
		newBytes := make([]byte, len(p.Get()))
		for i, b := range p.Get() {
			newBytes[len(newBytes)-i-1] = b
		}
		p.Set(newBytes)
		return nil
	})
	return []types.Message{newMsg}, nil
}

// CloseAsync shuts down the processor and stops processing requests.
func (m *Reverse) CloseAsync() {
}

// WaitForClose blocks until the processor has closed down.
func (m *Reverse) WaitForClose(timeout time.Duration) error {
	return nil
}

//------------------------------------------------------------------------------
