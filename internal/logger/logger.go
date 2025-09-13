package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string, service string, version string, env string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		logger = zap.NewNop()
	}

	logger = logger.With(
		zap.String("service", service),
		zap.String("version", version),
		zap.String("env", env),
	)

	return logger
}

func NewDevelopment() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func SetGlobal(logger *zap.Logger) {
	zap.ReplaceGlobals(logger)
}

func Sync() {
	if logger := zap.L(); logger != nil {
		logger.Sync()
	}
}
