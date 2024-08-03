package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	SeriesType string

	Series struct {
		seriesType SeriesType
		subType    string
		operation  string
		status     string
	}

	seriesContextKey struct{}
)

func (st SeriesType) String() string {
	return string(st)
}

const (
	SeriesTypeRPCHandler      SeriesType = "rpc_handler"
	SeriesTypeApiHandler      SeriesType = "api_handler"
	SeriesTypeUseCase         SeriesType = "use_case"
	SeriesTypeClient          SeriesType = "client"
	SeriesTypeDB              SeriesType = "postgres"
	SeriesTypeDatabusConsumer SeriesType = "databus_consumer"
)

// NewSeries creates a new Series instance with the given type and name.
func NewSeries(st SeriesType, subType string) Series {
	return Series{
		seriesType: st,
		subType:    subType,
		operation:  "undefined",
	}
}

// FromContext retrieves the Series from the context.
func FromContext(ctx context.Context) Series {
	series, ok := ctx.Value(seriesContextKey{}).(Series)
	if !ok {
		return Series{}
	}

	return series
}

// WithOperation sets the operation name in the Series and returns an updated context.
func (s Series) WithOperation(ctx context.Context, operation string) (context.Context, Series) {
	series := FromContext(ctx)

	if s.seriesType == series.seriesType &&
		s.subType == series.subType {
		series = Series{
			seriesType: s.seriesType,
			subType:    s.subType,
			operation:  series.appendOperation(operation),
		}

		return series.ToContext(ctx), series
	}

	series = Series{
		seriesType: s.seriesType,
		subType:    s.subType,
		operation:  operation,
	}

	return series.ToContext(ctx), series
}

// ToContext adds the Series to the context.
func (s Series) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, seriesContextKey{}, s)
}

const (
	seriesTypeInfo    = "info"
	seriesTypeSuccess = "success"
	seriesTypeError   = "error"
)

// Info returns the metric name and labels for an informational event.
func (s Series) Info() (string, prometheus.Labels) {
	return "operation_count", prometheus.Labels{
		"series_type": s.seriesType.String(),
		"sub_type":    s.subType,
		"operation":   s.operation,
		"status":      seriesTypeInfo,
	}
}

// Success returns the metric name and labels for a success event.
func (s Series) Success() (string, prometheus.Labels) {
	return "operation_count", prometheus.Labels{
		"series_type": s.seriesType.String(),
		"sub_type":    s.subType,
		"operation":   s.operation,
		"status":      seriesTypeSuccess,
	}
}

// Error returns the metric name and labels for an error event.
func (s Series) Error(message string) (string, prometheus.Labels) {
	return "operation_count", prometheus.Labels{
		"series_type":   s.seriesType.String(),
		"sub_type":      s.subType,
		"operation":     s.operation,
		"status":        seriesTypeError,
		"error_message": message,
	}
}

// Duration returns the metric name and labels for recording a duration.
func (s Series) Duration() (string, prometheus.Labels) {
	return "operation_duration_seconds", prometheus.Labels{
		"series_type": s.seriesType.String(),
		"sub_type":    s.subType,
		"operation":   s.operation,
	}
}

func (s Series) appendOperation(operation string) string {
	return s.operation + "_" + operation
}
