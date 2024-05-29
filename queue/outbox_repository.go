package queue

import (
	"context"
	"fmt"

	"github.com/gateway-fm/scriptorium/transactions"
)

// OutboxEvent represents an event stored in the outbox table.
type OutboxEvent struct {
	ID        int       `pg:",pk"`        // Primary key
	Data      []byte    `pg:"data"`       // Data column
	Topic     string    `pg:"topic"`      // Topic column
	Retry     int       `pg:"retry"`      // Retry column
	NextRetry uint      `pg:"next_retry"` // NextRetry column as minutes
	AckStatus AckStatus `pg:"ack_status"` // AckStatus column
	CreatedAt int64     `pg:"created_at"` // CreatedAt column as Unix timestamp
	UpdatedAt int64     `pg:"updated_at"` // UpdatedAt column as Unix timestamp
}

// OutboxRepository provides methods to interact with the outbox table.
type OutboxRepository struct {
	transactionFactory transactions.TransactionFactory
}

// NewOutboxRepository creates a new OutboxRepository.
func NewOutboxRepository(transactionFactory transactions.TransactionFactory) *OutboxRepository {
	return &OutboxRepository{transactionFactory: transactionFactory}
}

// InsertEvent inserts a new event into the outbox table or updates it if it already exists.
func (r *OutboxRepository) InsertEvent(ctx context.Context, event *OutboxEvent) error {
	_, err := r.transactionFactory.Transaction(ctx).
		Model(event).
		OnConflict("(data, topic) DO UPDATE").
		Set("retry = EXCLUDED.retry, next_retry = EXCLUDED.next_retry, ack_status = EXCLUDED.ack_status, updated_at = EXCLUDED.updated_at").
		Returning("*").
		Insert()
	if err != nil {
		return fmt.Errorf("insert event into outbox: %w", err)
	}
	return nil
}

// LoadPendingEvents loads all pending events from the outbox table.
func (r *OutboxRepository) LoadPendingEvents(ctx context.Context) ([]*OutboxEvent, error) {
	var events []*OutboxEvent
	query := r.transactionFactory.Transaction(ctx).
		Model(&events).
		Where("ack_status = ?", NACK)

	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("load pending events from outbox: %w", err)
	}
	return events, nil
}

// UpdateEventStatus updates the status and retry count of an event in the outbox table.
func (r *OutboxRepository) UpdateEventStatus(ctx context.Context, event *OutboxEvent) error {
	_, err := r.transactionFactory.Transaction(ctx).
		Model(event).
		Column("retry", "next_retry", "ack_status").
		Where("id = ?", event.ID).
		Update()
	if err != nil {
		return fmt.Errorf("update event status in outbox: %w", err)
	}
	return nil
}

// MarkEventAsProcessed marks an event as processed in the outbox table.
func (r *OutboxRepository) MarkEventAsProcessed(ctx context.Context, eventID int) error {
	_, err := r.transactionFactory.Transaction(ctx).
		Model(&OutboxEvent{}).
		Set("ack_status = ?", ACK).
		Where("id = ?", eventID).
		Update()
	if err != nil {
		return fmt.Errorf("mark event as processed in outbox: %w", err)
	}
	return nil
}

// MarkEventAsFailed marks an event as failed in the outbox table.
func (r *OutboxRepository) MarkEventAsFailed(ctx context.Context, eventID int) error {
	_, err := r.transactionFactory.Transaction(ctx).
		Model(&OutboxEvent{}).
		Set("ack_status = ?", NACK).
		Where("id = ?", eventID).
		Update()
	if err != nil {
		return fmt.Errorf("mark event as failed in outbox: %w", err)
	}
	return nil
}
