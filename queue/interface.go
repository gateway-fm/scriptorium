package queue

import (
	"context"
	"time"

	"github.com/gateway-fm/scriptorium/transactions"

	"github.com/gateway-fm/scriptorium/clog"
)

// EventBus defines an interface for subscribing to topics, publishing events, and managing event processing.
type EventBus interface {
	Subscribe(topic string, handler EventHandler, delays []int, durationType time.Duration)
	Publish(topic string, data []byte)
	StartProcessing(ctx context.Context) error
	Stop()
	ExceededMaxRetries(event *Event) bool
	SetLogger(log *clog.CustomLogger)
	AddEventToCtx(ctx context.Context, event *Event) context.Context
	WithOutbox(factory transactions.TransactionFactory)
}
