package metrics

import (
	"context"
	"fmt"
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

func (s Series) Success() string {
	return fmt.Sprintf("%s_%s_%s_success", s.st.String(), s.name, s.operation)
}

func (s Series) Error(errCode string) string {
	return fmt.Sprintf("%s.%s.%s.error.%s", s.st.String(), s.name, s.operation, errCode)
}

func (s Series) Duration() string {
	return fmt.Sprintf("%s.%s.%s.duration", s.st.String(), s.name, s.operation)
}

func (s Series) Info(code string) string {
	return fmt.Sprintf("%s.%s.%s.info.%s", s.st.String(), s.name, s.operation, code)
}

func (s Series) Operation() string {
	return s.operation
}

func (s Series) appendOperation(operation string) string {
	return s.operation + "." + operation
}
