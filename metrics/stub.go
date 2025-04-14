package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// registryStub is a no-op implementation of the Registry interface.
type registryStub struct{}

// NewRegistryStub creates a new instance of a no-op registry.
func NewRegistryStub() Registry {
	return &registryStub{}
}

// Inc is a no-op increment method for the stub registry.
func (s *registryStub) Inc(_ string, _ prometheus.Labels) {}

// RecordDuration is a no-op method for recording durations in the stub registry.
func (s *registryStub) RecordDuration(_ string, _ prometheus.Labels, _ float64) {}

// PrometheusRegistry returns nil for the stub registry.
func (s *registryStub) PrometheusRegistry() *prometheus.Registry {
	return nil
}
