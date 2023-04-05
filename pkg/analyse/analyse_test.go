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

type mockAnalyzer struct {
	mock.Mock
}

func (m *mockAnalyzer) Analyze(ctx context.Context, ecosystem string, manifestPath string, json, verbose bool) error {
	args := m.Called(ctx, ecosystem, manifestPath, json, verbose)
	return args.Error(0)
}

func TestStackReport(t *testing.T) {
	viper.Set(config.KeyConsentTelemetry.ToString(), false)
	t.Run("when analyzer fails should return an error", func(t *testing.T) {
		ctx := telemetry.GetContext(context.Background())

		// mock the analyzer
		analyzer := new(mockAnalyzer)
		analyzer.On("Analyze", ctx, "testecosystem", "fake-path", false, false).Return(errors.New("this is a fake error"))

		// create a fake manifest stubbed with the mocked analyzer
		manifest := Manifest{"fake.filename", "testecosystem", analyzer}

		err := StackReport(ctx, &manifest, "fake-path", false, false)
		require.Error(t, err)
		assert.Equal(t, "this is a fake error", err.Error())
		analyzer.AssertExpectations(t)
	})
}
