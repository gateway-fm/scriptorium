package metrics

import (
	"net/http"
	"sync/atomic"

	"github.com/gateway-fm/scriptorium/clog"
)

type HealthChecker struct {
	isReady   atomic.Value
	isHealthy atomic.Value
	logger    clog.CLog
}

func NewHealthChecker(logger clog.CLog) *HealthChecker {
	hc := &HealthChecker{
		logger: logger,
	}
	hc.isReady.Store(false)
	hc.isHealthy.Store(true)
	return hc
}

func (hc *HealthChecker) SetReady(ready bool) {
	hc.isReady.Store(ready)
}

func (hc *HealthChecker) SetHealthy(healthy bool) {
	hc.isHealthy.Store(healthy)
}

func (hc *HealthChecker) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	if hc.isHealthy.Load().(bool) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			hc.logger.ErrorCtx(r.Context(), err, "Failed to write liveness response")
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("unhealthy"))
		if err != nil {
			hc.logger.ErrorCtx(r.Context(), err, "Failed to write liveness response")
		}
	}
}

func (hc *HealthChecker) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	if hc.isReady.Load().(bool) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			hc.logger.ErrorCtx(r.Context(), err, "Failed to write readiness response")
		}
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err := w.Write([]byte("not ready"))
		if err != nil {
			hc.logger.ErrorCtx(r.Context(), err, "Failed to write readiness response")
		}
	}
}
