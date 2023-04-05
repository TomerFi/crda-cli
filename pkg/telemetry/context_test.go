package telemetry

import (
	"context"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func init() {
	utils.ConfigureLogging(false)
}

func TestTelemetryContext(t *testing.T) {
	t.Run("set and get a bool telemetry property", func(t *testing.T) {
		var ctx context.Context
		require.NotPanics(t, func() { ctx = GetContext(context.Background()) })
		require.NotPanics(t, func() { SetProperty(ctx, "testBoolVal", true) })

		val, ok := GetProperty(ctx, "testBoolVal")
		assert.True(t, ok)

		boolVal, err := strconv.ParseBool(val)
		assert.NoError(t, err)
		assert.True(t, boolVal)
	})

	t.Run("set and get an int telemetry property", func(t *testing.T) {
		var ctx context.Context
		require.NotPanics(t, func() { ctx = GetContext(context.Background()) })
		require.NotPanics(t, func() { SetProperty(ctx, "testIntVal", 1234) })

		val, ok := GetProperty(ctx, "testIntVal")
		assert.True(t, ok)

		intVal, err := strconv.Atoi(val)
		assert.NoError(t, err)
		assert.Equal(t, 1234, intVal)
	})

	t.Run("set and get a string telemetry property", func(t *testing.T) {
		var ctx context.Context
		require.NotPanics(t, func() { ctx = GetContext(context.Background()) })
		require.NotPanics(t, func() { SetProperty(ctx, "testIntVal", "a_string") })

		val, ok := GetProperty(ctx, "testIntVal")
		assert.True(t, ok)
		assert.Equal(t, "a_string", val)
	})
}

func TestPropertyKeyType_ToString(t *testing.T) {
	assert.Equal(t, string(KeyClient), KeyClient.ToString())
}
