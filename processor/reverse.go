package processor

import (
	"bytes"
	"context"

	"github.com/benthosdev/benthos/v4/public/service"
)

func init() {
	// Config spec is empty for now as we don't have any dynamic fields.
	configSpec := service.NewConfigSpec()

	constructor := func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
		return newReverseProcessor(mgr.Logger(), mgr.Metrics()), nil
	}

	err := service.RegisterProcessor("reverse", configSpec, constructor)
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type reverseProcessor struct {
	logger           *service.Logger
	countPalindromes *service.MetricCounter
}

func newReverseProcessor(logger *service.Logger, metrics *service.Metrics) *reverseProcessor {
	// The logger and metrics components will already be labelled with the
	// identifier of this component within a config.
	return &reverseProcessor{
		logger:           logger,
		countPalindromes: metrics.NewCounter("palindromes"),
	}
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
		r.countPalindromes.Incr(1)
	}

	m.SetBytes(newBytes)
	return []*service.Message{m}, nil
}

func (r *reverseProcessor) Close(ctx context.Context) error {
	return nil
}
