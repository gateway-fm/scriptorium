package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	SeriesType string

	Series struct {
		seriesType SeriesType
		subType    string
		operation  string
		status     string
		labels     prometheus.Labels
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

// NewSeries creates a new Series instance with the given type and subType.
func NewSeries(st SeriesType, subType string) Series {
	return Series{
		seriesType: st,
		subType:    subType,
		operation:  "undefined",
		labels:     make(prometheus.Labels),
	}
}

// WithLabels adds custom labels to the Series.
func (s Series) WithLabels(labels prometheus.Labels) Series {
	for k, v := range labels {
		s.labels[k] = v
	}
	return s
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
			labels:     series.labels,
		}

		return series.ToContext(ctx), series
	}

	series = Series{
		seriesType: s.seriesType,
		subType:    s.subType,
		operation:  operation,
		labels:     s.labels,
	}

	return series.ToContext(ctx), series
}

// ToContext adds the Series to the context.
func (s Series) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, seriesContextKey{}, s)
}

const (
	seriesTypeInfo     = "info"
	seriesTypeSuccess  = "success"
	seriesTypeError    = "error"
	seriesTypeDuration = "duration"
)

// Info returns the metric name and labels for an informational event.
func (s Series) Info(message string) (string, prometheus.Labels) {
	labels := prometheus.Labels{
		"series_type": s.seriesType.String(),
		"sub_type":    s.subType,
		"operation":   s.operation,
		"status":      seriesTypeInfo,
		"message":     message,
	}
	return "operation_count", mergeLabels(labels, s.labels)
}

// Success returns the metric name and labels for a success event.
func (s Series) Success() (string, prometheus.Labels) {
	labels := prometheus.Labels{
		"series_type": s.seriesType.String(),
		"sub_type":    s.subType,
		"operation":   s.operation,
		"status":      seriesTypeSuccess,
		"message":     "",
	}
	return "operation_count", mergeLabels(labels, s.labels)
}

// Error returns the metric name and labels for an error event.
func (s Series) Error(message string) (string, prometheus.Labels) {
	labels := prometheus.Labels{
		"series_type": s.seriesType.String(),
		"sub_type":    s.subType,
		"operation":   s.operation,
		"status":      seriesTypeError,
		"message":     message,
	}
	return "operation_count", mergeLabels(labels, s.labels)
}

// Duration returns the metric name and labels for recording a duration, where msg carries any additional information
func (s Series) Duration(d time.Duration, msg string) (string, prometheus.Labels, float64) {
	labels := prometheus.Labels{
		"series_type": s.seriesType.String(),
		"sub_type":    s.subType,
		"operation":   s.operation,
		"status":      seriesTypeDuration,
		"message":     msg,
	}

	return "operation_duration_seconds", mergeLabels(labels, s.labels), d.Seconds()
}

// appendOperation appends the operation to the Series operation string.
func (s Series) appendOperation(operation string) string {
	return s.operation + "_" + operation
}

// mergeLabels merges a set of additional labels into the base labels.
func mergeLabels(base prometheus.Labels, additional prometheus.Labels) prometheus.Labels {
	for k, v := range additional {
		base[k] = v
	}
	return base
}
