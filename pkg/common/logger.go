package common

import "go.uber.org/zap"

func InitLogger() *zap.Logger {
	logger, err := zap.NewProduction()

	if err != nil {
		panic("Logger could not initialized")
	}

	zap.ReplaceGlobals(logger)
	return logger
}
