package log

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	jsonEncoding = "json"
)

var Logger *zap.Logger
var once sync.Once

func NewLogger(logLevel zapcore.Level) {
	once.Do(func() {
		loggingConfig := zap.NewProductionConfig()
		loggingConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		loggingConfig.Encoding = jsonEncoding
		loggingConfig.Level.SetLevel(logLevel)
		logger, _ := loggingConfig.Build()
		zap.ReplaceGlobals(logger)
		Logger = logger
	})
}
