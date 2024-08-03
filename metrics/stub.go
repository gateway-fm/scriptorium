package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type registryStub struct{}

func NewRegistryStub() Registry {
	return &registryStub{}
}

func (s *registryStub) Inc(_ Series, _ string) {}

func (s *registryStub) RecordDuration(_ Series, _ float64) {
}

func (s *registryStub) PrometheusRegistry() *prometheus.Registry {
	return nil
}
