package input

import (
	"context"
	"math/rand"

	"github.com/Jeffail/benthos/v3/public/service"
)

func init() {
	configSpec := service.NewConfigSpec().
		Summary("Creates an input that generates garbage.").
		Field(service.NewIntField("length").Default(100))

	constructor := func(conf *service.ParsedConfig, mgr *service.Resources) (service.Input, error) {
		length, err := conf.FieldInt("length")
		if err != nil {
			return nil, err
		}
		return service.AutoRetryNacks(&gibberishInput{length}), nil
	}

	err := service.RegisterInput("gibberish", configSpec, constructor)
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type gibberishInput struct {
	length int
}

func (g *gibberishInput) Connect(ctx context.Context) error {
	return nil
}

func (g *gibberishInput) Read(ctx context.Context) (*service.Message, service.AckFunc, error) {
	b := make([]byte, g.length)
	for k := range b {
		b[k] = byte((rand.Int() % 94) + 32)
	}
	return service.NewMessage(b), func(ctx context.Context, err error) error {
		// Nacks are retries automatically when we use service.AutoRetryNacks
		return nil
	}, nil
}

func (g *gibberishInput) Close(ctx context.Context) error {
	return nil
}
