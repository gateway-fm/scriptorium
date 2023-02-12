package logger

import (
	"context"
	"encoding/json"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gateway-fm/scriptorium/helper"
)

const MaxAttributeChars = 200

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

// Log is invoking Zap Logger function
func Log() *Zaplog {
	initLogger()
	return instance
}

// LogWithContext is invoking Zap Logger function with context
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

func (z *Zaplog) InfoRedact(s string, fields ...zapcore.Field) *Zaplog {
	z.Info(z.Redact(s), fields...)
	return z
}

func (z *Zaplog) DebugRedact(s string, fields ...zapcore.Field) *Zaplog {
	z.Debug(z.Redact(s), fields...)
	return z
}

func (z *Zaplog) WarnRedact(s string, fields ...zapcore.Field) *Zaplog {
	z.Warn(z.Redact(s), fields...)
	return z
}

func (z *Zaplog) ErrorRedact(s string, fields ...zapcore.Field) *Zaplog {
	z.Error(z.Redact(s), fields...)
	return z
}

func (z *Zaplog) PanicRedact(s string, fields ...zapcore.Field) *Zaplog {
	z.Panic(z.Redact(s), fields...)
	return z
}

func (z *Zaplog) FatalRedact(s string, fields ...zapcore.Field) *Zaplog {
	z.Fatal(z.Redact(s), fields...)
	return z
}

func (z *Zaplog) InfoErr(e error, fields ...zapcore.Field) *Zaplog {
	z.Logger.Info(e.Error(), fields...)
	return z
}

func (z *Zaplog) DebugErr(e error, fields ...zapcore.Field) *Zaplog {
	z.Logger.Debug(e.Error(), fields...)
	return z
}

func (z *Zaplog) WarnErr(e error, fields ...zapcore.Field) *Zaplog {
	z.Logger.Warn(e.Error(), fields...)
	return z
}

func (z *Zaplog) ErrorErr(e error, fields ...zapcore.Field) *Zaplog {
	z.Logger.Error(e.Error(), fields...)
	return z
}

func (z *Zaplog) PanicErr(e error, fields ...zapcore.Field) *Zaplog {
	z.Logger.Panic(e.Error(), fields...)
	return z
}

// AnyCropped takes a key and an arbitrary value and chooses the best way to represent them as a field.
// It crops the content if it is larger than 200 characters.
func AnyCropped(key string, value interface{}) zap.Field {
	valueJson, err := json.Marshal(value)
	if err != nil || len(valueJson) <= MaxAttributeChars {
		return zap.Any(key, value)
	}
	return zap.ByteString(key, valueJson[0:MaxAttributeChars])
}
