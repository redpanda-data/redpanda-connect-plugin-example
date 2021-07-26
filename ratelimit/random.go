package ratelimit

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Jeffail/benthos/v3/public/service"
)

func init() {
	type randomRLConfig struct {
		MaxDuration string `yaml:"maximum_duration"`
	}

	configSpec := service.NewStructConfigSpec(func() interface{} {
		return &randomRLConfig{
			MaxDuration: "1s",
		}
	})

	constructor := func(conf *service.ParsedConfig, mgr *service.Resources) (service.RateLimit, error) {
		c := conf.AsStruct().(*randomRLConfig)
		maxDuration, err := time.ParseDuration(c.MaxDuration)
		if err != nil {
			return nil, fmt.Errorf("invalid max duration: %w", err)
		}
		return &randomRateLimit{maxDuration}, nil
	}

	err := service.RegisterRateLimit("random", configSpec, constructor)
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type randomRateLimit struct {
	max time.Duration
}

func (r *randomRateLimit) Access(context.Context) (time.Duration, error) {
	return time.Duration(rand.Int() % int(r.max)), nil
}

func (r *randomRateLimit) Close(ctx context.Context) error {
	return nil
}
