package metrics

import "github.com/prometheus/client_golang/prometheus"

type registryStub struct{}

func NewRegistryStub() Registry {
	return &registryStub{}
}

func (s *registryStub) Inc(_ string) {}

func (s *registryStub) RecordDuration(_ string, _ []string) *prometheus.HistogramVec {
	return nil
}

func (s *registryStub) PrometheusRegistry() *prometheus.Registry {
	return nil
}
