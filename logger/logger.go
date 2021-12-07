package logger


import (
	"context"
	"ohdearcodingisfun/helper"
	"sync"

	"go.uber.org/zap"
)

type Zaplog struct {
	*zap.Logger
}

var instance *Zaplog
var once sync.Once

//Log is invoking Zap Logger function
func Log() *Zaplog {
	once.Do(func() {
		logger, _ := zap.NewProduction()
		instance = &Zaplog{logger}
	})
	return instance
}

//LogWithContext is invoking Zap Logger function with context
func LogWithContext(ctx context.Context) *zap.Logger {
	once.Do(func() {
		logger, _ := zap.NewProduction()
		instance = &Zaplog{logger}
	})
	return instance.With(zap.String(string(helper.ContextKeyRequestID), helper.GetRequestID(ctx)))
}
