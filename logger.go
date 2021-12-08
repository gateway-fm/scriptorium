package scriptorium

import (
	"fmt"
	"github.com/gateway-fm/scriptorium/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// ILogger implements logger's functions
type ILogger interface {
	InitLogger() error
}
type Logger struct {
	logger *logger.Zaplog
	Env    string
}

// Logger constructor
func NewLogger(logenv string) ILogger {
	return &Logger{
		logger: logger.Log(),
		Env:    logenv,
	}
}

// InitLogger initializes logger on
// application level
func (c *Logger) InitLogger() error {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level == zapcore.ErrorLevel
		}),
	)
	env, err := logger.EnvFromStr(c.Env)
	if err != nil {
		return fmt.Errorf("failed when getting the app env: %w", err)
	}
	switch env {
	case logger.Local:
		//experimental
		//cfg := zap.NewDevelopmentConfig()
		//cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		//cfg.EncoderConfig.TimeKey = ""
		//cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		//	logger.Log().Logger, _ = cfg.Build()
		loclog := zap.New(core)

		c.logger.Logger = loclog
	case logger.Production, logger.Development:
		c.logger.Logger, _ = zap.NewProduction()
	default:
		return fmt.Errorf("env wasn't specified correctly: %v", c.Env)
	}
	return nil
}

func (c *Logger) Suppress() {
	zap.AddCaller()
}
