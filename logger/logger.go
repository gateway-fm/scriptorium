package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gateway-fm/scriptorium/helper"
)

type Zaplog struct {
	*zap.Logger
}

var (
	instance *Zaplog
	once     sync.Once
	appEnv   AppEnv
)

// SetLoggerMode set Logger level from given string
func SetLoggerMode(envStr string) {
	appEnv = EnvFromStr(envStr)
}

//Log is invoking Zap Logger function
func Log() *Zaplog {
	initLogger()
	return instance
}

//LogWithContext is invoking Zap Logger function with context
func LogWithContext(ctx context.Context) *zap.Logger {
	initLogger()
	return instance.With(zap.String(string(helper.ContextKeyRequestID), helper.GetRequestID(ctx)))
}

// initLogger initialise Logger instance only once
func initLogger() {
	once.Do(func() {
		switch appEnv {
		case Local:
			cfg := zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			cfg.EncoderConfig.TimeKey = ""
			cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
			log, _ := cfg.Build()
			instance = &Zaplog{log}
		default:
			log, _ := zap.NewProduction()
			instance = &Zaplog{log}
		}
	})
}
