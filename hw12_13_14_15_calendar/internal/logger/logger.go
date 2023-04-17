package logger

import (
	"os"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InvalidLevel
	case "warn":
		return zapcore.WarnLevel
	default:
		return zapcore.ErrorLevel
	}
}

func New(level string) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	return zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg),
		os.Stdout, parseLevel(level)))
}
