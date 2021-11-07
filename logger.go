package logger

import (
	"go.uber.org/zap"
	"sync"
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
