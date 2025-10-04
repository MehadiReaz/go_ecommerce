package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

// Init initializes the logger
func Init() {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	log = logger.Sugar()
}

// Info logs an info message
func Info(msg string, keysAndValues ...interface{}) {
	log.Infow(msg, keysAndValues...)
}

// Debug logs a debug message
func Debug(msg string, keysAndValues ...interface{}) {
	log.Debugw(msg, keysAndValues...)
}

// Error logs an error message
func Error(msg string, keysAndValues ...interface{}) {
	log.Errorw(msg, keysAndValues...)
}

// Warn logs a warning message
func Warn(msg string, keysAndValues ...interface{}) {
	log.Warnw(msg, keysAndValues...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, keysAndValues ...interface{}) {
	log.Fatalw(msg, keysAndValues...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return log.Sync()
}
