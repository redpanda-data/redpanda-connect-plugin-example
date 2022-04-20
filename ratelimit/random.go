package ratelimit

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/benthosdev/benthos/v4/public/service"
)

func init() {
	configSpec := service.NewConfigSpec().
		Field(service.NewStringField("maximum_duration").Default("1s"))

	constructor := func(conf *service.ParsedConfig, mgr *service.Resources) (service.RateLimit, error) {
		durStr, err := conf.FieldString("maximum_duration")
		if err != nil {
			return nil, err
		}

		maxDuration, err := time.ParseDuration(durStr)
		if err != nil {
			return nil, fmt.Errorf("invalid max duration: %w", err)
		}

		return newRandomRateLimit(maxDuration), nil
	}

	if err := service.RegisterRateLimit("random", configSpec, constructor); err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type randomRateLimit struct {
	max time.Duration
}

func newRandomRateLimit(max time.Duration) *randomRateLimit {
	return &randomRateLimit{max}
}

func (r *randomRateLimit) Access(context.Context) (time.Duration, error) {
	return time.Duration(rand.Int() % int(r.max)), nil
}

func (r *randomRateLimit) Close(ctx context.Context) error {
	return nil
}
