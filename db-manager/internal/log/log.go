package log

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var once sync.Once

func NewLogger() {
	// Ensures this is only ran once per execution
	once.Do(func() {
		// Creates log with Debug level
		loggingConfig := zap.NewProductionConfig()
		loggingConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		loggingConfig.Level.SetLevel(zapcore.DebugLevel)

		// Specify the log file path for output
		logFilePath := "server/server.log"
		loggingConfig.OutputPaths = []string{logFilePath}

		// Create the log file and clear any information from the prior execution
		logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			// Handle error
			panic(err)
		}

		// Finishes setting up logger
		defer logFile.Close()
		logger, _ := loggingConfig.Build()
		zap.ReplaceGlobals(logger)
		Logger = logger
	})
}
