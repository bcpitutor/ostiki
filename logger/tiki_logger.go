package logger

import (
	"github.com/bcpitutor/ostiki/appconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TikiLogger struct {
	Logger *zap.Logger
}

func GetTikiLogger(cfg *appconfig.AppConfig) *TikiLogger {
	outputFile := cfg.LoggerOutputFile

	conf := zap.Config{
		Encoding: "json",
		Level:    zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths: []string{
			outputFile,
			"stdout",
		},
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "file",
			MessageKey:   "message",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	l, _ := conf.Build()

	return &TikiLogger{
		Logger: l,
	}
}
