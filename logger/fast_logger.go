package logger

import (
	"github.com/gateway-fm/scriptorium/fast_helper"
	"github.com/gateway-fm/scriptorium/helper"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// LogWithFastContext is invoking Zap Logger function with fasthttp context
func LogWithFastContext(ctx *fasthttp.RequestCtx) *Zaplog {
	initLogger()

	publicKey := fast_helper.GetPublicKey(ctx)
	if publicKey == "" {
		publicKey = helper.PublicKeyNotSet
	}

	return instance.With(
		zap.String(string(helper.ContextKeyRequestID), fast_helper.GetRequestID(ctx)),
		zap.String(string(helper.PublicKey), publicKey),
	)
}

func LogFast(ctx *fasthttp.RequestCtx) *Zaplog {
	initLogger()
	return instance
}
