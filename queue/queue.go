package queue

import (
	"context"
	"sync"
	"time"

	"github.com/gateway-fm/scriptorium/clog"
)

type AckStatus string

const (
	ACK  AckStatus = "ACK"  // message acknowledged, no need to retry.
	NACK AckStatus = "NACK" // message not acknowledged, need to retry.
)

// Event represents a message or event that can be published to a topic within the EventBus.
type Event struct {
	Data      []byte        // Data is the binary payload of the event.
	Retry     int           // Retry indicates how many times this event has been retried.
	Topic     string        // Topic is the name of the topic to which the event is published.
	NextRetry time.Duration // NextRetry specifies the delay before the next retry attempt.
	AckStatus AckStatus     // AckStatus specifies whether the message is done.
}

// EventHandler is a function type that processes an Event and returns an error if the processing fails.
type EventHandler func(ctx context.Context, event Event) AckStatus

// eventBus implements the EventBus interface with support for topic-based subscriptions and event retries.
type eventBus struct {
	ctx      context.Context            // ctx is the base context for all operations.
	cf       context.CancelFunc         // cf is a function to cancel the context, used for stopping the event processing.
	log      *clog.CustomLogger         // log is a custom logger for logging information about event processing.
	handlers map[string][]EventHandler  // handlers store a slice of event handlers for each topic.
	delay    map[string][]time.Duration // delay specifies the retry delays for each topic.
	queue    chan Event                 // queue is the channel through which events are published and processed.
	lock     sync.RWMutex               // lock is used to synchronize access to handlers and delays.
}

// NewEventBus creates a new instance of an eventBus with a specified buffer size for the event queue and attaches a logger.
// The context passed is used to manage the lifecycle of the event processing.
func NewEventBus(ctx context.Context, size int) EventBus {
	ctx, cf := context.WithCancel(ctx)
	return &eventBus{
		ctx:      ctx,
		cf:       cf,
		handlers: make(map[string][]EventHandler),
		delay:    make(map[string][]time.Duration),
		queue:    make(chan Event, size),
	}
}

func (bus *eventBus) SetLogger(log *clog.CustomLogger) {
	bus.log = log
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

// Publish sends an event with the specified data to the specified topic.
func (bus *eventBus) Publish(topic string, data []byte) {
	bus.queue <- Event{Data: data, Topic: topic, Retry: 0, NextRetry: 0}
}

// StartProcessing begins processing events from the queue. It listens for cancellation via the provided context to gracefully stop processing.
func (bus *eventBus) StartProcessing(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			bus.log.InfoContext(ctx, "Processing stopped due to context cancellation")
			return
		case event, ok := <-bus.queue:
			if !ok {
				bus.log.InfoContext(ctx, "Processing stopped, queue channel closed")
				return
			}
			go processEvent(ctx, bus, event)
		}
	}
}

// processEvent handles the processing of a single event, including retry logic and error handling.
func processEvent(ctx context.Context, bus *eventBus, event Event) {
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

			go bus.retryEvent(ctx, event)
		case event.Retry >= maxRetries:
			bus.log.DebugfCtx(ctx, "Max retries for event: %+v\n", event)
		case status == ACK:
			bus.log.DebugfCtx(ctx, "Message read: %+v\n", event)
		}
	}
}

// retryEvent attempts to re-enqueue an event for processing after a delay, respecting the provided context.
func (bus *eventBus) retryEvent(ctx context.Context, event Event) {
	select {
	case <-ctx.Done():
		bus.log.DebugfCtx(ctx, "Retry canceled due to context cancellation for event: %+v\n", event)
		return
	case <-time.After(event.NextRetry):
		select {
		case bus.queue <- event:
			bus.log.DebugfCtx(ctx, "Event re-enqueued after delay: %s, topic %s\n", event.Data, event.Topic)
		case <-ctx.Done():
			bus.log.DebugfCtx(ctx, "Failed to enqueue event due to context cancellation: %s, topic %s\n", event.Data, event.Topic)
		}
	}
}

// Stop triggers the stopping of the event bus processing by cancelling the context.
func (bus *eventBus) Stop() {
	bus.cf()
}

func (bus *eventBus) ReachedMaxRetries(event Event) bool {
	return event.Retry >= len(bus.delay[event.Topic])
}
