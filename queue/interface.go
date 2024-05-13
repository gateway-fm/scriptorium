package queue

import (
	"context"
	"time"

	"github.com/gateway-fm/scriptorium/clog"
)

// EventBus defines an interface for subscribing to topics, publishing events, and managing event processing.
type EventBus interface {
	Subscribe(topic string, handler EventHandler, delays []int, durationType time.Duration)
	Publish(topic string, data []byte)
	StartProcessing(ctx context.Context)
	Stop()
	ReachedMaxRetries(event Event) bool
	SetLogger(log *clog.CustomLogger)
}
