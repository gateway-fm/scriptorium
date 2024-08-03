package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type registry struct {
	Subsystem    string
	Namespace    string
	PromRegistry *prometheus.Registry

	metricsMu  sync.Mutex
	counters   map[string]prometheus.Counter
	histograms map[string]*prometheus.HistogramVec
}

func NewRegistry(subsystem, namespace string) Registry {
	r := &registry{
		Subsystem:    subsystem,
		Namespace:    namespace,
		PromRegistry: prometheus.NewRegistry(),
		counters:     make(map[string]prometheus.Counter),
		histograms:   make(map[string]*prometheus.HistogramVec),
	}

	registerMetrics(r)

	return r
}

func (r *registry) Inc(name string) {
	r.metricsMu.Lock()
	defer r.metricsMu.Unlock()

	counter, exists := r.counters[name]
	if !exists {
		counter = prometheus.NewCounter(prometheus.CounterOpts{
			Subsystem: r.Subsystem,
			Namespace: r.Namespace,
			Name:      name,
		})
		r.PromRegistry.MustRegister(counter)
		r.counters[name] = counter
	}
	counter.Inc()
}

func (r *registry) RecordDuration(name string, labels []string) *prometheus.HistogramVec {
	r.metricsMu.Lock()
	defer r.metricsMu.Unlock()

	histogram, exists := r.histograms[name]
	if !exists {
		histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: r.Subsystem,
			Namespace: r.Namespace,
			Name:      name,
			Buckets:   prometheus.DefBuckets,
		}, labels)
		r.PromRegistry.MustRegister(histogram)
		r.histograms[name] = histogram
	}
	return histogram
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
