package fast_helper

import (
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
type ContextKey string

// ContextKeyRequestID is the ContextKey for RequestID
const ContextKeyRequestID ContextKey = "requestID"
const PublicKey ContextKey = "userPublicKey"

const RequestIDPrefix string = "reqid://"

// SetRandomRequestID will attach a brand new request ID to a http request
func SetRandomRequestID(ctx *fasthttp.RequestCtx) {
	reqID, err := uuid.NewV4()
	if err != nil {
		ctx.Logger().Printf("couldn't create randomly generated UUID ")
	}
	ctx.SetUserValue(ContextKeyRequestID, RequestIDPrefix+reqID.String())
}

func SetRequestID(ctx *fasthttp.RequestCtx, requestID string) {
	ctx.SetUserValue(ContextKeyRequestID, requestID)
}

// SetResponseCode sets http response status code on the context
func SetResponseCode(ctx *fasthttp.RequestCtx, code int) {
	ctx.Response.SetStatusCode(code)
}

// GetResponseCode returns http response status code from the context
func GetResponseCode(ctx *fasthttp.RequestCtx) int {
	return ctx.Response.StatusCode()
}

// GetPublicKey returns public key from the context
func GetPublicKey(ctx *fasthttp.RequestCtx) string {
	publicKey := ctx.Value(PublicKey)
	if ret, ok := publicKey.(string); ok {
		return ret
	}
	return ""
}

// GetRequestID will get reqID from a http request and return it as a string
func GetRequestID(ctx *fasthttp.RequestCtx) string {
	reqID := ctx.Value(ContextKeyRequestID)
	if ret, ok := reqID.(string); ok {
		return ret
	}
	return ""
}

type DebugKey string

// Debug is the DebugKey for debug mode
const Debug DebugKey = "debug"

// SetDebug will set debug mode to a http request
func SetDebug(ctx *fasthttp.RequestCtx, debug bool) {
	ctx.SetUserValue(Debug, debug)
}

// GetDebug will get debug mode from a http request and return it
func GetDebug(ctx *fasthttp.RequestCtx) bool {
	debug := ctx.Value(Debug)
	if ret, ok := debug.(bool); ok {
		return ret
	}
	return false
}

func GetFeatTag() bool {
	switch viper.GetString("ENV") {
	case "dev":
		return true
	case "local":
		return true
	default:
		return false
	}
}
