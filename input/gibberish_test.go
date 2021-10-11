package input

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGibberishInputConfigValidation(t *testing.T) {
	// First element is a config, second element is a string expected in the
	// config error.
	for _, test := range [][2]string{
		{`length: 0`, `length must be greater than 0`},
		{`length: 100000`, `that length is way too high bruh`},
		{`length: 10`, ``},
		{``, ``}, // Using the default length is fine
	} {
		conf, err := gibberishConfigSpec.ParseYAML(test[0], nil)
		require.NoError(t, err)

		_, err = newGibberishInput(conf)
		if test[1] == "" {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			assert.Contains(t, err.Error(), test[1])
		}
	}
}

func TestGibberishInput(t *testing.T) {
	gibConf, err := gibberishConfigSpec.ParseYAML(`length: 10`, nil)
	require.NoError(t, err)

	gibInput, err := newGibberishInput(gibConf)
	require.NoError(t, err)

	require.NoError(t, gibInput.Connect(context.Background()))

	msg, ackFn, err := gibInput.Read(context.Background())
	require.NoError(t, err)

	msgBytes, err := msg.AsBytes()
	require.NoError(t, err)

	assert.Len(t, msgBytes, 10)
	require.NoError(t, ackFn(context.Background(), nil))

	assert.NoError(t, gibInput.Close(context.Background()))
}
