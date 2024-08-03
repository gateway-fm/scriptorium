package metrics

import "github.com/prometheus/client_golang/prometheus"

// Registry defines the interface for managing metrics.
type Registry interface {
	Inc(series Series, status string)
	RecordDuration(series Series, duration float64)
	PrometheusRegistry() *prometheus.Registry
}
