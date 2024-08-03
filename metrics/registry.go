package metrics

import (
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type registry struct {
	Subsystem    string
	Namespace    string
	PromRegistry *prometheus.Registry

	metricsMu  sync.Mutex
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
}

func NewRegistry(subsystem, namespace string) Registry {
	r := &registry{
		Subsystem:    subsystem,
		Namespace:    namespace,
		PromRegistry: prometheus.NewRegistry(),
		counters:     make(map[string]*prometheus.CounterVec),
		histograms:   make(map[string]*prometheus.HistogramVec),
	}

	registerMetrics(r)

	return r
}

func (r *registry) sanitizeMetricName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

func (r *registry) Inc(series Series, status string) {
	r.metricsMu.Lock()
	defer r.metricsMu.Unlock()

	metricName, labels := series.SuccessLabels()
	labels["status"] = status
	sanitized := r.sanitizeMetricName(metricName)
	counter, exists := r.counters[sanitized]
	if !exists {
		counter = prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem: r.Subsystem,
			Namespace: r.Namespace,
			Name:      sanitized,
		}, []string{"series_type", "name", "operation", "status"})
		r.PromRegistry.MustRegister(counter)
		r.counters[sanitized] = counter
	}
	counter.With(labels).Inc()
}

func (r *registry) RecordDuration(series Series, duration float64) {
	r.metricsMu.Lock()
	defer r.metricsMu.Unlock()

	metricName, labels := series.DurationLabels()
	sanitized := r.sanitizeMetricName(metricName)
	histogram, exists := r.histograms[sanitized]
	if !exists {
		histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: r.Subsystem,
			Namespace: r.Namespace,
			Name:      sanitized,
			Buckets:   prometheus.DefBuckets,
		}, []string{"series_type", "name", "operation"})
		r.PromRegistry.MustRegister(histogram)
		r.histograms[sanitized] = histogram
	}
	histogram.With(labels).Observe(duration)
}

func (r *registry) PrometheusRegistry() *prometheus.Registry {
	return r.PromRegistry
}

func registerMetrics(registry *registry) {
	registry.PromRegistry.MustRegister(
		collectors.NewGoCollector(
			collectors.WithGoCollectorMemStatsMetricsDisabled(),
			collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsScheduler),
		))
}
