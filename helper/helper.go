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

const RequestIDPrefix string = "reqid://"

// SetRequestID will attach a brand new request ID to a http request
func SetRequestID(ctx context.Context) context.Context {
	reqID, err := uuid.NewV4()
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, ContextKeyRequestID, RequestIDPrefix+reqID.String())
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

// DVFAuthentication is used for context.Context value.
type DVFAuthentication string

// DvfAuthToken is the DVFAuthentication indicates
const DvfAuthToken DVFAuthentication = "dvfAuthToken"

func GetDVFAuthToken(ctx context.Context) string {
	dvfAuthToken := ctx.Value(DvfAuthToken)
	if token, ok := dvfAuthToken.(string); ok {
		return token
	}
	return ""
}
