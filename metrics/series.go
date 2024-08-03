package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	SeriesType string

	Series struct {
		st        SeriesType
		name      string
		operation string
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

func NewSeries(st SeriesType, name string) Series {
	return Series{
		st:        st,
		name:      name,
		operation: "undefined",
	}
}

func FromContext(ctx context.Context) Series {
	series, ok := ctx.Value(seriesContextKey{}).(Series)
	if !ok {
		return Series{}
	}

	return series
}

func (s Series) WithOperation(ctx context.Context, operation string) (context.Context, Series) {
	series := FromContext(ctx)

	if s.st == series.st &&
		s.name == series.name {
		series = Series{
			st:        s.st,
			name:      s.name,
			operation: series.appendOperation(operation),
		}

		return series.ToContext(ctx), series
	}

	series = Series{
		st:        s.st,
		name:      s.name,
		operation: operation,
	}

	return series.ToContext(ctx), series
}

func (s Series) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, seriesContextKey{}, s)
}

func (s Series) SuccessLabels() (string, prometheus.Labels) {
	return "success_count", prometheus.Labels{
		"series_type": s.st.String(),
		"name":        s.name,
		"operation":   s.operation,
		"status":      "success",
	}
}

func (s Series) ErrorLabels(errCode string) (string, prometheus.Labels) {
	return "error_count", prometheus.Labels{
		"series_type": s.st.String(),
		"name":        s.name,
		"operation":   s.operation,
		"status":      "error",
		"error_code":  errCode,
	}
}

func (s Series) DurationLabels() (string, prometheus.Labels) {
	return "operation_duration_seconds", prometheus.Labels{
		"series_type": s.st.String(),
		"name":        s.name,
		"operation":   s.operation,
	}
}

func (s Series) InfoLabels(code string) (string, prometheus.Labels) {
	return "info_events", prometheus.Labels{
		"series_type": s.st.String(),
		"name":        s.name,
		"operation":   s.operation,
		"info_code":   code,
	}
}

func (s Series) Operation() string {
	return s.operation
}

func (s Series) appendOperation(operation string) string {
	return s.operation + "." + operation
}
