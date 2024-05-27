package queue

import (
	"context"
	"sync"
	"time"

	"github.com/gateway-fm/scriptorium/transactions"

	"github.com/gateway-fm/scriptorium/clog"
)

type AckStatus string

const (
	ACK  AckStatus = "ACK"  // message acknowledged, no need to retry.
	NACK AckStatus = "NACK" // message not acknowledged, need to retry.
)

// Event represents a message or event that can be published to a topic within the EventBus.
type Event struct {
	ID        int           // ID is the identifier for the event in the database.
	Data      []byte        // Data is the binary payload of the event.
	Retry     int           // Retry indicates how many times this event has been retried.
	Topic     string        // Topic is the name of the topic to which the event is published.
	NextRetry time.Duration // NextRetry specifies the delay before the next retry attempt.
	AckStatus AckStatus     // AckStatus specifies whether the message is done.
}

// EventHandler is a function type that processes an Event and returns an error if the processing fails.
type EventHandler func(ctx context.Context, event *Event) AckStatus

// eventBus implements the EventBus interface with support for topic-based subscriptions and event retries.
type eventBus struct {
	ctx      context.Context            // ctx is the base context for all operations.
	cf       context.CancelFunc         // cf is a function to cancel the context, used for stopping the event processing.
	log      *clog.CustomLogger         // log is a custom logger for logging information about event processing.
	handlers map[string][]EventHandler  // handlers store a slice of event handlers for each topic.
	delay    map[string][]time.Duration // delay specifies the retry delays for each topic.
	queue    chan *Event                // queue is the channel through which events are published and processed.
	lock     sync.RWMutex               // lock is used to synchronize access to handlers and delays.
	outbox   *OutboxRepository
}

// NewEventBus creates a new instance of an eventBus with a specified buffer size for the event queue and attaches a logger.
func NewEventBus(ctx context.Context, size int) EventBus {
	ctx, cf := context.WithCancel(ctx)
	return &eventBus{
		ctx:      ctx,
		cf:       cf,
		handlers: make(map[string][]EventHandler),
		delay:    make(map[string][]time.Duration),
		queue:    make(chan *Event, size),
	}
}

func (bus *eventBus) SetLogger(log *clog.CustomLogger) {
	bus.log = log
}

func (bus *eventBus) WithOutbox(factory transactions.TransactionFactory) {
	bus.outbox = NewOutboxRepository(factory)
}

// Subscribe adds an event handler for a specific topic with predefined retry delays.
func (bus *eventBus) Subscribe(
	topic string,
	handler EventHandler,
	delays []int,
	durationType time.Duration,
) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	delaysDuration := make([]time.Duration, len(delays))

	for index, delay := range delays {
		delaysDuration[index] = time.Duration(delay) * durationType
	}

	bus.handlers[topic] = append(bus.handlers[topic], handler)
	bus.delay[topic] = delaysDuration
}

func (bus *eventBus) Publish(topic string, data []byte) {
	event := &Event{
		Data:      data,
		Topic:     topic,
		Retry:     0,
		NextRetry: 0,
		AckStatus: NACK,
	}

	if bus.outbox != nil {
		outboxEvent := convertEventToOutboxEvent(event)
		if err := bus.outbox.InsertEvent(bus.ctx, outboxEvent); err != nil {
			bus.log.Error("Failed to save event to outbox: %v", err)
			return
		}
		event.ID = outboxEvent.ID
	}

	bus.queue <- event
}

// StartProcessing begins processing events from the queue. It listens for cancellation via the provided context to gracefully stop processing.
func (bus *eventBus) StartProcessing(ctx context.Context) error {
	err := bus.loadEventsFromOutbox(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			bus.log.InfoContext(ctx, "Processing stopped due to context cancellation")
			return nil
		case event, ok := <-bus.queue:
			if !ok {
				bus.log.InfoContext(ctx, "Processing stopped, queue channel closed")
				return nil
			}
			go processEvent(ctx, bus, event)
		}
	}
}

// processEvent handles the processing of a single event, including retry logic and error handling.
func processEvent(ctx context.Context, bus *eventBus, event *Event) {
	ctx = bus.AddEventToCtx(ctx, event)

	handlers, ok := bus.handlers[event.Topic]
	if !ok {
		return
	}
	for _, handler := range handlers {
		status := handler(ctx, event)
		maxRetries := len(bus.delay[event.Topic])
		switch {
		case status == NACK && event.Retry < maxRetries:
			event.Retry++
			event.NextRetry = bus.delay[event.Topic][event.Retry-1]
			if bus.outbox != nil {
				bus.updateEventStatus(ctx, event)
			}
			go bus.retryEvent(ctx, event)
		case event.Retry >= maxRetries:
			bus.log.DebugCtx(ctx, "Max retries for event")
			if bus.outbox != nil {
				bus.markEventAsFailed(ctx, event.ID)
			}
		case status == ACK:
			bus.log.DebugCtx(ctx, "Message acknowledged")
			if bus.outbox != nil {
				bus.markEventAsProcessed(ctx, event.ID)
			}
		}
	}
}

// retryEvent attempts to re-enqueue an event for processing after a delay, respecting the provided context.
func (bus *eventBus) retryEvent(ctx context.Context, event *Event) {
	select {
	case <-ctx.Done():
		bus.log.DebugCtx(ctx, "Retry canceled due to context cancellation for event: %+v\n", event)
		return
	case <-time.After(event.NextRetry):
		select {
		case bus.queue <- event:
			bus.log.DebugCtx(ctx, "Event re-enqueued after delay")
		case <-ctx.Done():
			bus.log.DebugCtx(ctx, "Failed to enqueue event due to context cancellation")
		}
	}
}

// Stop triggers the stopping of the event bus processing by cancelling the context.
func (bus *eventBus) Stop() {
	bus.cf()
}

func (bus *eventBus) ExceededMaxRetries(event *Event) bool {
	return event.Retry > len(bus.delay[event.Topic])
}

func (bus *eventBus) AddEventToCtx(ctx context.Context, event *Event) context.Context {
	return bus.log.AddKeysValuesToCtx(ctx, map[string]interface{}{
		"event_data":              string(event.Data),
		"event_retry":             event.Retry,
		"event_topic":             event.Topic,
		"event_nextRetry_minutes": event.NextRetry.Minutes(),
		"event_ackStatus":         event.AckStatus,
	})
}

// convertEventToOutboxEvent converts an Event to an OutboxEvent.
func convertEventToOutboxEvent(event *Event) *OutboxEvent {
	return &OutboxEvent{
		ID:        event.ID,
		Data:      event.Data,
		Topic:     event.Topic,
		Retry:     event.Retry,
		NextRetry: uint(event.NextRetry.Minutes()),
		AckStatus: event.AckStatus,
	}
}

// convertOutboxEventToEvent converts an OutboxEvent to an Event.
func convertOutboxEventToEvent(outboxEvent *OutboxEvent) *Event {
	return &Event{
		ID:        outboxEvent.ID,
		Data:      outboxEvent.Data,
		Topic:     outboxEvent.Topic,
		Retry:     outboxEvent.Retry,
		NextRetry: time.Duration(outboxEvent.NextRetry) * time.Minute,
		AckStatus: outboxEvent.AckStatus,
	}
}
