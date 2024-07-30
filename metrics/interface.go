package metrics

import "github.com/prometheus/client_golang/prometheus"

type Registry interface {
	Inc(name string)
	RecordDuration(name string, labels []string) *prometheus.HistogramVec
	PrometheusRegistry() *prometheus.Registry
}
