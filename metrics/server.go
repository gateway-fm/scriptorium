package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/gateway-fm/scriptorium/clog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsEndpoint   = "/metrics"
	livenessEndpoint  = "/healthz"
	readinessEndpoint = "/readyz"
)

type Server interface {
	Start(ctx context.Context, address string)
	Stop(ctx context.Context)
}

type server struct {
	logger      clog.CLog
	registry    Registry
	healthCheck *HealthChecker
	srv         *http.Server
}

func NewServer(logger clog.CLog, registry Registry, healthCheck *HealthChecker) *server {
	return &server{
		logger:      logger,
		registry:    registry,
		healthCheck: healthCheck,
	}
}

func (s *server) start(ctx context.Context, address string) {
	mux := http.NewServeMux()

	mux.Handle(metricsEndpoint, promhttp.HandlerFor(s.registry.PrometheusRegistry(), promhttp.HandlerOpts{}))
	mux.HandleFunc(livenessEndpoint, s.healthCheck.LivenessHandler)
	mux.HandleFunc(readinessEndpoint, s.healthCheck.ReadinessHandler)

	ctx = s.logger.AddKeysValuesToCtx(ctx, map[string]interface{}{
		"metrics_address": address,
	})

	s.srv = &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go func() {
		s.logger.InfoCtx(ctx, "Metrics and Health Check server is running")

		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.ErrorCtx(ctx, err, "Failed to start metrics server")
		}
	}()
}

func (s *server) Start(ctx context.Context, address string) {
	go s.start(ctx, address)
}

func (s *server) Stop(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.ErrorCtx(ctx, err, "could not shutdown the metrics server")
	}
}
