package ratelimit

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Jeffail/benthos/v3/lib/log"
	"github.com/Jeffail/benthos/v3/lib/metrics"
	"github.com/Jeffail/benthos/v3/lib/ratelimit"
	"github.com/Jeffail/benthos/v3/lib/types"
)

//------------------------------------------------------------------------------

func init() {
	ratelimit.RegisterPlugin(
		"random",
		func() interface{} {
			conf := NewRandomConfig()
			return &conf
		},
		func(
			iconf interface{},
			mgr types.Manager,
			logger log.Modular,
			stats metrics.Type,
		) (types.RateLimit, error) {
			conf, ok := iconf.(*RandomConfig)
			if !ok {
				return nil, errors.New("failed to cast config")
			}
			return NewRandom(*conf, mgr, logger, stats)
		},
	)
	ratelimit.DocumentPlugin(
		"random",
		`Randomly throttles by a specified duration based on a random number
generator.`,
		nil,
	)
}

//------------------------------------------------------------------------------

// RandomConfig contains config fields for the random rate limit.
type RandomConfig struct {
	ThrottleFor string `json:"throttle_for" yaml:"throttle_for"`
}

// NewRandomConfig returns a RandomConfig with default values.
func NewRandomConfig() RandomConfig {
	return RandomConfig{
		ThrottleFor: "1s",
	}
}

// Random is a rate limit that randomly throttles messages.
type Random struct {
	throttleFor time.Duration
}

// NewRandom returns a Random cache.
func NewRandom(
	conf RandomConfig, mgr types.Manager, log log.Modular, stats metrics.Type,
) (types.RateLimit, error) {
	tFor, err := time.ParseDuration(conf.ThrottleFor)
	if err != nil {
		return nil, err
	}
	return &Random{
		throttleFor: tFor,
	}, nil
}

//------------------------------------------------------------------------------

// Access the rate limited resource. Returns a duration or an error if the rate
// limit check fails. The returned duration is either zero (meaning the resource
// can be accessed) or a reasonable length of time to wait before requesting
// again.
func (r *Random) Access() (time.Duration, error) {
	if rand.Int()%2 == 0 {
		return r.throttleFor, nil
	}
	return 0, nil
}

// CloseAsync shuts down the rate limit.
func (r *Random) CloseAsync() {
}

// WaitForClose blocks until the rate limit has closed down.
func (r *Random) WaitForClose(timeout time.Duration) error {
	return nil
}

//------------------------------------------------------------------------------
