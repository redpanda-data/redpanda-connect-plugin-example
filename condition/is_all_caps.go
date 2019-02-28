package condition

import (
	"bytes"
	"errors"

	"github.com/Jeffail/benthos/lib/log"
	"github.com/Jeffail/benthos/lib/metrics"
	"github.com/Jeffail/benthos/lib/processor/condition"
	"github.com/Jeffail/benthos/lib/types"
)

//------------------------------------------------------------------------------

func init() {
	condition.RegisterPlugin(
		"is_all_caps",
		func() interface{} {
			conf := NewIsAllCapsConfig()
			return &conf
		},
		func(
			iconf interface{},
			mgr types.Manager,
			logger log.Modular,
			stats metrics.Type,
		) (types.Condition, error) {
			conf, ok := iconf.(*IsAllCapsConfig)
			if !ok {
				return nil, errors.New("failed to cast config")
			}
			return NewIsAllCaps(*conf, mgr, logger, stats)
		},
	)
	condition.DocumentPlugin(
		"is_all_caps",
		`Checks whether a message of a batch is all capital letters.`,
		nil,
	)
}

//------------------------------------------------------------------------------

// IsAllCapsConfig is a configuration struct containing fields for the IsAllCaps
// condition.
type IsAllCapsConfig struct {
	Part int `json:"part" yaml:"part"`
}

// NewIsAllCapsConfig returns a IsAllCapsConfig with default values.
func NewIsAllCapsConfig() IsAllCapsConfig {
	return IsAllCapsConfig{
		Part: 0,
	}
}

//------------------------------------------------------------------------------

// IsAllCaps is a condition that checks whether a message part is all capital
// letters.
type IsAllCaps struct {
	stats metrics.Type
	log   log.Modular
	part  int
}

// NewIsAllCaps returns a IsAllCaps condition.
func NewIsAllCaps(
	conf IsAllCapsConfig, mgr types.Manager, log log.Modular, stats metrics.Type,
) (condition.Type, error) {
	return &IsAllCaps{
		stats: stats,
		log:   log,
		part:  conf.Part,
	}, nil
}

//------------------------------------------------------------------------------

// Check attempts to check a message part against a configured condition.
func (c *IsAllCaps) Check(msg types.Message) bool {
	return bytes.Equal(bytes.ToUpper(msg.Get(c.part).Get()), msg.Get(c.part).Get())
}

//------------------------------------------------------------------------------
