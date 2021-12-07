package logger

import (
	"fmt"
	"github.com/gateway-fm/logger/logger"
	"go.uber.org/zap"
)

type IConvey interface {
	InitLogger() error
}
type Convey struct {
	logger *logger.Zaplog
	Env    string
}

func NewConvey(logenv string) IConvey {
	convey := &Convey{
		logger: logger.Log(),
		Env:    logenv,
	}

	return convey

}

func (c *Convey) InitLogger() error {
	env, err := logger.EnvFromStr(c.Env)
	if err != nil {
		return fmt.Errorf("failed when getting the app env: %w", err)
	}
	switch env {
	case logger.Local:
		c.logger.Logger, _ = zap.NewDevelopment()
	case logger.Production, logger.Development:
		c.logger.Logger, _ = zap.NewProduction()
	default:
		return fmt.Errorf("env wasn't specified correctly: %v", c.Env)
	}
	return nil
}
