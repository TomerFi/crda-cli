package utils

import "go.uber.org/zap"

var Logger *zap.SugaredLogger

// ConfigureLogging is used to configure the global logger with or without debugging
func ConfigureLogging(debug bool) {
	// configure logging
	logConfig := zap.NewProductionConfig()
	logConfig.OutputPaths = []string{"stdout"}
	if debug {
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, _ := logConfig.Build()
	defer logger.Sync() // //nolint:errcheck
	Logger = logger.Sugar()
}
