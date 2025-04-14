package tracing

import "time"

type TraceConfig struct {
	ServiceName string
	Env         string

	Discovery *Discovery
	OtlpRetry *OtlpRetry
}

type Discovery struct {
	Driver      string
	ConsulAddr  string
	Manual      []string
	Transport   string
	BackendName string
	Opts        *DiscoveryOptions
}

type DiscoveryOptions struct {
	IsOptional   bool
	OptionalPath string
}

type OtlpRetry struct {
	Enabled         bool
	InitialInterval time.Duration
	MaxInterval     time.Duration
	MaxElapsedTime  time.Duration
}
