package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initialises the global Zap logger.
// Call once from main before any other package uses it.
func Init(production bool) error {
	var cfg zap.Config
	if production {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	var err error
	log, err = cfg.Build()
	return err
}

// Get returns the global logger.
func Get() *zap.Logger {
	if log == nil {
		// Fallback: never panic in library code.
		log, _ = zap.NewDevelopment()
	}
	return log
}

// Sync flushes any buffered log entries.
func Sync() { _ = log.Sync() }
