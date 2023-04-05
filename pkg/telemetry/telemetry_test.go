package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"

	"github.com/segmentio/analytics-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSegmentClient struct {
	mock.Mock
}

func (m *mockSegmentClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockSegmentClient) Enqueue(message analytics.Message) error {
	m.Called(message)
	return nil
}

func init() {
	utils.ConfigureLogging(false)
}

func TestPushEvent(t *testing.T) {
	t.Run("with consent and no segment client should not return an error", func(t *testing.T) {
		ctx := GetContext(context.Background())
		viper.Set(config.KeyConsentTelemetry.ToString(), true)
		assert.NoError(t, PushEvent(ctx, nil, "aa12-bb34", "dummy-event1", time.Now()))
	})

	t.Run("without consent should not push events", func(t *testing.T) {
		ctx := GetContext(context.Background())
		viper.Set(config.KeyConsentTelemetry.ToString(), false)

		mockSegment := new(mockSegmentClient)
		mockSegment.On("Close").Return(nil) // only stub Close, Enqueue will not be invoked

		require.NoError(t, PushEvent(ctx, mockSegment, "cc56-dd78", "dummy-event2", time.Now()))

		mockSegment.AssertNotCalled(t, "Enqueue")
		mockSegment.AssertExpectations(t)
	})

	t.Run("with consent and a client should push identify and track events", func(t *testing.T) {
		ctx := GetContext(context.Background())
		viper.Set(config.KeyConsentTelemetry.ToString(), true)

		mockSegment := new(mockSegmentClient)
		mockSegment.On("Close").Return(nil) // stub the Close method

		// stub for the first Enqueue invocation identifying the user
		identifyEventMsg := mock.MatchedBy(func(message analytics.Identify) bool { return message.UserId == "ee90-ff12" })
		mockSegment.On("Enqueue", identifyEventMsg).Return(nil).Once()

		// stub for the second Enqueue invocation pushing the track event
		trackEventMsg := mock.MatchedBy(func(message analytics.Track) bool {
			return message.UserId == "ee90-ff12" &&
				message.Event == "dummy-event3"
		})
		mockSegment.On("Enqueue", trackEventMsg).Return(nil).Once()

		require.NoError(t, PushEvent(ctx, mockSegment, "ee90-ff12", "dummy-event3", time.Now()))

		mockSegment.AssertExpectations(t)
	})

}
