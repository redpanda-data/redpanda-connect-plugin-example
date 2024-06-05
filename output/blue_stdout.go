package output

import (
	"context"
	"fmt"

	"github.com/redpanda-data/benthos/v4/public/service"
)

func init() {
	err := service.RegisterOutput(
		"blue_stdout", service.NewConfigSpec(),
		func(conf *service.ParsedConfig, mgr *service.Resources) (out service.Output, maxInFlight int, err error) {
			return &blueOutput{}, 1, nil
		})
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type blueOutput struct{}

func (b *blueOutput) Connect(ctx context.Context) error {
	return nil
}

func (b *blueOutput) Write(ctx context.Context, msg *service.Message) error {
	content, err := msg.AsBytes()
	if err != nil {
		return err
	}
	fmt.Printf("\033[01;34m%s\033[m\n", content)
	return nil
}

func (b *blueOutput) Close(ctx context.Context) error {
	return nil
}
