package tracing

import (
	"context"
	"fmt"
)

var defaultTracerName string

func SetDefaultTracerName(name string) {
	defaultTracerName = name
}

func GetDefaultTracer() ITracer {
	return getOpenTelemetryTracer(defaultTracerName)
}

func GetTracer(name string) ITracer {
	return getOpenTelemetryTracer(name)
}

func NewTracerProvider(conf *TraceConfig) (ITracerProvider, error) {
	if conf == nil {
		return nil, fmt.Errorf("the trace config is nil while trying to get a new tracer provider")
	}
	return getOpenTelemetryTracerProvider(conf)
}

func ContextWithSpanContext(ctx context.Context, traceID string) (context.Context, error) {
	return contextWithSpanContext(ctx, traceID)
}

func GetTraceIDFromContext(ctx context.Context) string {
	return getTraceIDFromContext(ctx)
}
