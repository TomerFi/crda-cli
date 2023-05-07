package analyse

import (
	"context"
	"errors"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/telemetry"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	utils.ConfigureLogging(false)
}

type mockProvider struct {
	mock.Mock
}

func (m *mockProvider) Provide(ctx context.Context, manifestPath string) ([]byte, string, error) {
	args := m.Called(ctx, manifestPath)
	return nil, "", args.Error(0)
}

func TestGetStackReport(t *testing.T) {
	t.Skip("WIP")
	viper.Set(config.KeyConsentTelemetry.ToString(), false)
	t.Run("when provider fails should return an error", func(t *testing.T) {
		ctx := telemetry.GetContext(context.Background())

		// mock the provider
		analyzer := new(mockProvider)
		analyzer.On("Provide", ctx, "fake-path").Return(errors.New("this is a fake error"))

		// create a fake manifest stubbed with the mocked analyzer
		manifest := Manifest{"fake.filename", "testecosystem", analyzer}

		err := GetStackReport(ctx, &manifest, "fake-path", false, false)
		require.Error(t, err)
		assert.Equal(t, "this is a fake error", err.Error())
		analyzer.AssertExpectations(t)
	})
}
