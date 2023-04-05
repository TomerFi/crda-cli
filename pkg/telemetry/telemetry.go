package telemetry

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/rhecosystemappeng/crda-cli/pkg/utils"

	"github.com/segmentio/analytics-go"
)

// PushEvent is used to push a track event
// the event properties will be fetched from the context
// the duration wil be calculated based on the start time
// returns error if failed identifying or pushing the track event to segment
func PushEvent(ctx context.Context, segmentClient analytics.Client, userId, eventName string, startTime time.Time) error {
	utils.Logger.Debug("pushing telemetry track event")
	if segmentClient != nil {
		defer segmentClient.Close() // close segment client when done
	}

	// if telemetry consent is false or not yet set, return without pushing telemetry
	if !isTelemetryConsent() {
		utils.Logger.Debug("pushing telemetry skipped, no consent given")
		return nil
	}

	// prepare event properties
	eventProps := analytics.NewProperties()
	// include collected event properties
	for k, v := range ctx.Value(contextKey).(contextValue).props {
		eventProps.Set(k.ToString(), v)
	}

	// add essential event properties
	eventProps.Set(KeyPlatform.ToString(), runtime.GOOS)
	eventProps.Set(KeyVersion.ToString(), utils.GetCRDAVersion())
	eventProps.Set(KeyDuration.ToString(), time.Since(startTime))

	if segmentClient != nil {
		// identify the user id with segment
		if err := segmentClient.Enqueue(analytics.Identify{UserId: userId}); err != nil {
			return fmt.Errorf("failed identifying the user id with segment, %w", err)
		}

		// push track event to segment
		if err := segmentClient.Enqueue(analytics.Track{
			UserId:     userId,
			Event:      eventName,
			Properties: eventProps,
		}); err != nil {
			return fmt.Errorf("failed to push telemetry to segment, %w", err)
		}
	}

	return nil
}
