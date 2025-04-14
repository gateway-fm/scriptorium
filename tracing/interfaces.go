package tracing

import (
	"context"
)

type ITracer interface {
	Start(ctx context.Context, spanName string) (context.Context, ISpan)
}

type ISpan interface {
	End()
	AddEvent(name string)
	SetStatus(code Code, description string)
	RecordError(err error)
	GetTraceID() string

	SetAttributeInt(k string, v int)
	SetAttributeString(k string, v string)
	SetAttributeBool(k string, v bool)
	SetAttributeFloat64(k string, v float64)
	SetAttributeStringSlice(k string, v []string)
}

type ITracerProvider interface {
	Shutdown(ctx context.Context) error
}
