package queue

import "context"

// loadEventsFromOutbox loads events from the outbox table into the in-memory queue.
func (bus *eventBus) loadEventsFromOutbox(ctx context.Context) error {
	if bus.outbox == nil {
		return nil
	}

	events, err := bus.outbox.LoadPendingEvents(ctx)
	if err != nil {
		bus.log.ErrorCtx(ctx, err, "Failed to load events from outbox")
		return err
	}

	for _, outboxEvent := range events {
		bus.queue <- convertOutboxEventToEvent(outboxEvent)
	}

	return nil
}

// updateEventStatus updates the status and retry count of an event in the outbox table.
func (bus *eventBus) updateEventStatus(ctx context.Context, event *Event) {
	outboxEvent := convertEventToOutboxEvent(event)

	if err := bus.outbox.UpdateEventStatus(ctx, outboxEvent); err != nil {
		bus.log.ErrorCtx(ctx, err, "Failed to update event status in outbox")
	}
}

// markEventAsProcessed marks an event as processed in the outbox table.
func (bus *eventBus) markEventAsProcessed(ctx context.Context, eventID int) {
	if err := bus.outbox.MarkEventAsProcessed(ctx, eventID); err != nil {
		bus.log.ErrorCtx(ctx, err, "Failed to mark event as processed in outbox")
	}
}

// markEventAsFailed marks an event as failed in the outbox table.
func (bus *eventBus) markEventAsFailed(ctx context.Context, eventID int) {
	if err := bus.outbox.MarkEventAsFailed(ctx, eventID); err != nil {
		bus.log.ErrorCtx(ctx, err, "Failed to mark event as failed in outbox")
	}
}
