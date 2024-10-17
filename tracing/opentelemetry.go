package tracing

import (
	"context"
	"fmt"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/gateway-fm/scriptorium/logger"
	"github.com/gateway-fm/service-pool/discovery"
)

func getOpenTelemetryTracer(tracerName string) *openTelemetryTracer {
	return &openTelemetryTracer{tracerName: tracerName}
}

type openTelemetryTracer struct {
	tracerName string
}

func (t *openTelemetryTracer) Start(ctx context.Context, spanName string) (context.Context, ISpan) {
	ctx, span := otel.Tracer(t.tracerName).Start(ctx, spanName)
	return ctx, &openTelemetrySpan{span: span}
}

type openTelemetrySpan struct {
	span trace.Span
}

func (s *openTelemetrySpan) End() {
	if s == nil || s.span == nil {
		return
	}
	s.span.End()
}

func (s *openTelemetrySpan) SetAttributeInt(k string, v int) {
	if s == nil || s.span == nil {
		return
	}
	s.span.SetAttributes(attribute.Int(k, v))
}

func (s *openTelemetrySpan) SetAttributeString(k string, v string) {
	if s == nil || s.span == nil {
		return
	}
	s.span.SetAttributes(attribute.String(k, v))
}

func (s *openTelemetrySpan) SetAttributeBool(k string, v bool) {
	if s == nil || s.span == nil {
		return
	}
	s.span.SetAttributes(attribute.Bool(k, v))
}

func (s *openTelemetrySpan) SetAttributeFloat64(k string, v float64) {
	if s == nil || s.span == nil {
		return
	}
	s.span.SetAttributes(attribute.Float64(k, v))
}

func (s *openTelemetrySpan) SetAttributeStringSlice(k string, v []string) {
	if s == nil || s.span == nil {
		return
	}
	s.span.SetAttributes(attribute.StringSlice(k, v))
}

func (s *openTelemetrySpan) AddEvent(name string) {
	if s == nil || s.span == nil {
		return
	}
	s.span.AddEvent(name)
}

func (s *openTelemetrySpan) SetStatus(code Code, description string) {
	if s == nil || s.span == nil {
		return
	}
	s.span.SetStatus(otelcodes.Code(code), description)
}

func (s *openTelemetrySpan) RecordError(err error) {
	if s == nil || s.span == nil {
		return
	}
	s.span.RecordError(err)
}

func (s *openTelemetrySpan) GetTraceID() string {
	if s == nil || s.span == nil {
		return ""
	}

	// .SpanContext() returns struct (not interface and not a pointer), so it cannot be nil
	return s.span.SpanContext().TraceID().String()
}

type openTelemetryTracerProvider struct {
	provider *sdktrace.TracerProvider
}

func (p *openTelemetryTracerProvider) Shutdown(ctx context.Context) error {
	if p.provider == nil {
		return fmt.Errorf("the open telemetry trace provider is nil while trying to Shutdown")
	}
	return p.provider.Shutdown(ctx)
}

func getOpenTelemetryTracerProvider(conf *TraceConfig) (*openTelemetryTracerProvider, error) {
	if conf == nil {
		return nil, fmt.Errorf("the trace config is nil while trying to get an open telemetry trace provider")
	}

	exporter, err := newTracerProviderExporter(conf)
	if err != nil {
		return nil, fmt.Errorf("an error occured while trying to create a new trace provider exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newTracerProviderResource(conf)))

	p := &openTelemetryTracerProvider{provider: tp}
	otel.SetTracerProvider(tp)

	return p, nil
}

func newTracerProviderExporter(conf *TraceConfig) (*otlptrace.Exporter, error) {
	disc, err := newTracerProviderDiscovery(conf)
	if err != nil {
		return nil, fmt.Errorf("an error occured while trying to get a new trace provider discovery: %w", err)
	}

	services, err := disc.Discover(conf.Discovery.BackendName)
	if err != nil {
		return nil, fmt.Errorf("an error occured while trying to discover distributed tracing backends: %w", err)
	}

	if len(services) <= 0 {
		return nil, fmt.Errorf("no active distributed tracing backends were discovered")
	}

	endpoint, _ := url.Parse(services[0].Address())

	clientAttributes := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint.Host),
		otlptracehttp.WithRetry(
			otlptracehttp.RetryConfig{
				Enabled:         conf.OtlpRetry.Enabled,
				InitialInterval: conf.OtlpRetry.InitialInterval,
				MaxInterval:     conf.OtlpRetry.MaxInterval,
				MaxElapsedTime:  conf.OtlpRetry.MaxElapsedTime,
			}),
	}
	if discovery.TransportFromString(conf.Discovery.Transport) == discovery.TransportHttp {
		clientAttributes = append(clientAttributes, otlptracehttp.WithInsecure())
	}

	client := otlptracehttp.NewClient(clientAttributes...)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}
	return exporter, nil
}

func newTracerProviderResource(conf *TraceConfig) *resource.Resource {
	if conf == nil {
		logger.Log().Error("the trace config is nil while trying to get a trace provider resource")
		return nil
	}

	attributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String(conf.ServiceName),
		attribute.String("env", conf.Env),
	}

	res, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			attributes...))
	return res
}

func newTracerProviderDiscovery(conf *TraceConfig) (discovery.IServiceDiscovery, error) {
	driver := discovery.DriverFromString(conf.Discovery.Driver)
	newDiscovery, err := discovery.ParseDiscoveryDriver(driver)
	if err != nil {
		return nil, fmt.Errorf("parse discovery driver: %w", err)
	}

	switch driver {
	case discovery.DriverConsul:
		return newDiscovery(discovery.TransportFromString(conf.Discovery.Transport), nil, conf.Discovery.ConsulAddr)
	case discovery.DriverManual:
		return newDiscovery(discovery.TransportFromString(conf.Discovery.Transport), nil, conf.Discovery.Manual...)
	}

	return nil, fmt.Errorf("the unknown discovery driver occured while trying to init discovery for a tracing provider")
}

func spanContextWithTraceID(traceID string) (*trace.SpanContext, error) {
	traceId, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return nil, fmt.Errorf("an error occured while trying to get TraceIDFromHex: %w", err)
	}

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceId})
	return &spanContext, nil
}

func contextWithSpanContext(ctx context.Context, traceID string) (context.Context, error) {
	spanContext, err := spanContextWithTraceID(traceID)
	if err != nil {
		return nil, fmt.Errorf("an error occured while trying to get a new context with SpanContext: %w", err)
	}

	return trace.ContextWithSpanContext(ctx, *spanContext), nil
}

func getTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		return span.SpanContext().TraceID().String()
	}
	return ""
}
