package logging

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	componentField = "component"
	logLevelEnv    = "LOG_LEVEL"
)

// NewLogger creates a new logger.
func NewLogger() (*zap.Logger, error) {
	logLevel := zapcore.InfoLevel
	logLevelString := os.Getenv(logLevelEnv)
	if logLevelString != "" {
		atomicLevel, err := zap.ParseAtomicLevel(logLevelString)
		if err != nil {
			return nil, fmt.Errorf("error parsing level")
		}
		logLevel = atomicLevel.Level()
	}
	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(logLevel)

	log, err := logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("error creating new logger: %w", err)
	}

	return log, nil
}

// WithComponent will return the logger with a component field attached.
func WithComponent(log *zap.Logger, name string) *zap.Logger {
	return log.With(zap.String(componentField, name))
}
