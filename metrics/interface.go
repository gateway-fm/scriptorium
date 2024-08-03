package metrics

import "github.com/prometheus/client_golang/prometheus"

type Registry interface {
	Inc(series Series, status string)
	RecordDuration(series Series, duration float64)
	PrometheusRegistry() *prometheus.Registry
}
