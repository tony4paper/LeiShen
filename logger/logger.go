package logger

import "go.uber.org/zap"

var (
	logger *zap.Logger
)

func init() {
	logger, _ = zap.NewProduction()
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}
