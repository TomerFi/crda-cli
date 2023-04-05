package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestConfigureLogging(t *testing.T) {
	t.Run("configuring for debugging should enable debug level", func(t *testing.T) {
		require.NotPanics(t, func() { ConfigureLogging(true) })
		assert.True(t, Logger.Desugar().Core().Enabled(zap.DebugLevel))
	})

	t.Run("configuring without debugging should not enable debug level", func(t *testing.T) {
		require.NotPanics(t, func() { ConfigureLogging(false) })
		assert.False(t, Logger.Desugar().Core().Enabled(zap.DebugLevel))
	})
}
