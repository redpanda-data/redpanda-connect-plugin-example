package processor

import (
	"bytes"
	"context"

	"github.com/Jeffail/benthos/v3/public/service"
)

func init() {
	constructor := func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
		return &reverseProcessor{
			logger: mgr.Logger(),
		}, nil
	}

	err := service.RegisterProcessor("reverse", service.NewConfigSpec(), constructor)
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type reverseProcessor struct {
	logger *service.Logger
}

func (r *reverseProcessor) Process(ctx context.Context, m *service.Message) (service.MessageBatch, error) {
	bytesContent, err := m.AsBytes()
	if err != nil {
		return nil, err
	}

	newBytes := make([]byte, len(bytesContent))
	for i, b := range bytesContent {
		newBytes[len(newBytes)-i-1] = b
	}

	if bytes.Equal(newBytes, bytesContent) {
		r.logger.Infof("Woah! This is like totally a palindrome: %s", bytesContent)
	}

	m.SetBytes(newBytes)
	return []*service.Message{m}, nil
}

func (r *reverseProcessor) Close(ctx context.Context) error {
	return nil
}
