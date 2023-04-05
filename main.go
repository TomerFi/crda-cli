package main

import (
	"github.com/rhecosystemappeng/crda-cli/cmd"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/telemetry"
	"github.com/segmentio/analytics-go"
	"os"
)

// TODO get this out of here, perhaps to an env var
const writeKey = "MW6rAYP7Q6AAiSAZ3Ussk6eMebbVcchD" // test

func main() {
	// we only enqueue 2 messages, one for identifying and one for tracking
	// the segment client is being checked for nil before usage
	segmentClient, _ := analytics.NewWithConfig(writeKey, analytics.Config{BatchSize: 2})
	userIdFile := telemetry.GetUserIdFilePath()        // $HOME/.redhat/anonymousId
	configDirectory := config.GetConfigDirectoryPath() // $HOME/.crda

	os.Exit(cmd.Run(segmentClient, userIdFile, configDirectory))
}
