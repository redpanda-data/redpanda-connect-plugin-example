package processor

import (
	"context"
	"testing"

	"github.com/benthosdev/benthos/v4/public/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReverseProcessor(t *testing.T) {
	// It's safe to pass nil in place of a logger for testing purposes
	revProc := newReverseProcessor(nil, nil)

	result, err := revProc.Process(context.Background(), service.NewMessage([]byte("hello world")))
	require.NoError(t, err)
	require.Len(t, result, 1)

	resBytes, err := result[0].AsBytes()
	require.NoError(t, err)
	assert.Equal(t, "dlrow olleh", string(resBytes))

	// Try a palindrome
	result, err = revProc.Process(context.Background(), service.NewMessage([]byte("wooooow")))
	require.NoError(t, err)
	require.Len(t, result, 1)

	resBytes, err = result[0].AsBytes()
	require.NoError(t, err)
	assert.Equal(t, "wooooow", string(resBytes))
}
