package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(env string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.StacktraceKey = "stacktrace"

	var err error
	Logger, err = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		logFatal("Failed to initialize Zap logger: %v", err)
	}

	// Ensure all buffered logs are flushed when the application exits.
	// This is important for graceful shutdown.
	zap.ReplaceGlobals(Logger) // Set as global Zap logger
	zap.RedirectStdLog(Logger) // Redirect standard library logs to Zap
	Logger.Info("Logger initialized successfully", zap.String("environment", env))
}

// logFatal is a helper to log a fatal error before the main logger is fully initialized.
func logFatal(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(1)
}
