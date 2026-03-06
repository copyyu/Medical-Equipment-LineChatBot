package event

import "context"

// EventBus defines the interface for publishing and subscribing to domain events.
// This interface lives in the domain layer so that usecases can depend on it
// without knowing about Redis or any specific message broker implementation.
type EventBus interface {
	// Publish sends an event to all subscribers
	Publish(ctx context.Context, event Event) error

	// Subscribe returns a channel that receives events matching the given types.
	// If no types are specified, all events are received.
	// The caller must cancel the context to unsubscribe and clean up resources.
	Subscribe(ctx context.Context, eventTypes ...string) (<-chan Event, error)

	// Close gracefully shuts down the event bus
	Close() error
}
