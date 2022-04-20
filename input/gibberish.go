package input

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/benthosdev/benthos/v4/public/service"
)

var gibberishConfigSpec = service.NewConfigSpec().
	Summary("Creates an input that generates garbage.").
	Field(service.NewIntField("length").Default(100))

func newGibberishInput(conf *service.ParsedConfig) (service.Input, error) {
	length, err := conf.FieldInt("length")
	if err != nil {
		return nil, err
	}
	if length <= 0 {
		return nil, fmt.Errorf("length must be greater than 0, got: %v", length)
	}
	if length > 10000 {
		return nil, errors.New("that length is way too high bruh")

	}
	return service.AutoRetryNacks(&gibberishInput{length}), nil
}

func init() {
	err := service.RegisterInput(
		"gibberish", gibberishConfigSpec,
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.Input, error) {
			return newGibberishInput(conf)
		})
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
		// Nacks are retried automatically when we use service.AutoRetryNacks
		return nil
	}, nil
}

func (g *gibberishInput) Close(ctx context.Context) error {
	return nil
}
