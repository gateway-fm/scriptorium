package helper

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
type ContextKey string

// ContextKeyRequestID is the ContextKey for RequestID
const ContextKeyRequestID ContextKey = "requestID"
const PublicKey ContextKey = "userPublicKey"

const RequestIDPrefix string = "reqid://"

// SetRandomRequestID will attach a brand new request ID to a http request
func SetRandomRequestID(ctx context.Context) context.Context {
	reqID, err := uuid.NewV4()
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, ContextKeyRequestID, RequestIDPrefix+reqID.String())
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

//func GetPublicKeyFromScriptorium(method func() []byte) []byte{
//	return method()
//}

// GetPublicKey
func GetPublicKey(ctx context.Context) string {
	publicKey := ctx.Value(PublicKey)
	if ret, ok := publicKey.(string); ok {
		return ret
	}
	return ""
}

// GetRequestID will get reqID from a http request and return it as a string
func GetRequestID(ctx context.Context) string {
	reqID := ctx.Value(ContextKeyRequestID)
	if ret, ok := reqID.(string); ok {
		return ret
	}
	return ""
}

func GetContextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	requestId := GetRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	ctx = context.WithValue(ctx, ContextKeyRequestID, requestId)
	return ctx, cancel
}

type DebugKey string

// Debug is the DebugKey for debug mode
const Debug DebugKey = "debug"

// SetDebug will set debug mode to a http request
func SetDebug(ctx context.Context, debug bool) context.Context {
	return context.WithValue(ctx, Debug, debug)
}

// GetDebug will get debug mode from a http request and return it
func GetDebug(ctx context.Context) bool {
	debug := ctx.Value(Debug)
	if ret, ok := debug.(bool); ok {
		return ret
	}
	return false
}
