package metrics

import "github.com/prometheus/client_golang/prometheus"

// Registry defines the interface for managing metrics.
type Registry interface {
	Inc(string, prometheus.Labels)
	RecordDuration(string, prometheus.Labels, float64)
	PrometheusRegistry() *prometheus.Registry
}
