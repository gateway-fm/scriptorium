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
	*redactor
}

var (
	instance *Zaplog
	once     sync.Once
	appEnv   AppEnv
)

// initLogger initialise Logger instance only once
func initLogger() {
	once.Do(func() {
		r := NewRedactor()
		switch appEnv {
		case Local:
			cfg := zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			cfg.EncoderConfig.TimeKey = ""
			cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
			log, _ := cfg.Build()
			instance = &Zaplog{log, r}
		default:
			log, _ := zap.NewProduction()
			instance = &Zaplog{log, r}
		}
	})
}

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
func LogWithContext(ctx context.Context) *Zaplog {
	initLogger()

	publicKey := helper.GetPublicKey(ctx)
	if publicKey == "" {
		publicKey = helper.PublicKeyNotSet
	}

	return instance.With(
		zap.String(string(helper.ContextKeyRequestID), helper.GetRequestID(ctx)),
		zap.String(string(helper.PublicKey), publicKey),
	)
}

func (z *Zaplog) With(fields ...zapcore.Field) *Zaplog {
	return &Zaplog{z.Logger.With(fields...), z.redactor}
}

func (z *Zaplog) InfoRedact(s string) *Zaplog {
	z.Info(z.Redact(s))
	return z
}

func (z *Zaplog) DebugRedact(s string) *Zaplog {
	z.Debug(z.Redact(s))
	return z
}

func (z *Zaplog) WarnRedact(s string) *Zaplog {
	z.Warn(z.Redact(s))
	return z
}

func (z *Zaplog) ErrorRedact(s string) *Zaplog {
	z.Error(z.Redact(s))
	return z
}

func (z *Zaplog) PanicRedact(s string) *Zaplog {
	z.Panic(z.Redact(s))
	return z
}

func (z *Zaplog) FatalRedact(s string) *Zaplog {
	z.Fatal(z.Redact(s))
	return z
}

func (z *Zaplog) InfoErr(e error) *Zaplog {
	z.Logger.Info(e.Error())
	return z
}

func (z *Zaplog) DebugErr(e error) *Zaplog {
	z.Logger.Debug(e.Error())
	return z
}

func (z *Zaplog) WarnErr(e error) *Zaplog {
	z.Logger.Warn(e.Error())
	return z
}

func (z *Zaplog) ErrorErr(e error) *Zaplog {
	z.Logger.Error(e.Error())
	return z
}

func (z *Zaplog) PanicErr(e error) *Zaplog {
	z.Logger.Panic(e.Error())
	return z
}
