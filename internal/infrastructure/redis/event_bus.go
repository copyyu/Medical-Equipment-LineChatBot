package redis

import (
	"context"
	"encoding/json"
	"log"
	"medical-webhook/internal/domain/event"
	"sync"

	goredis "github.com/redis/go-redis/v9"
)

const channelName = "medical:events"

// EventBus implements event.EventBus using Redis Pub/Sub
type EventBus struct {
	client *goredis.Client
	mu     sync.Mutex
	closed bool
}

// NewEventBus creates a new Redis-backed EventBus
func NewEventBus(client *goredis.Client) *EventBus {
	return &EventBus{
		client: client,
	}
}

// Publish sends an event to the Redis channel
func (eb *EventBus) Publish(ctx context.Context, e event.Event) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	if err := eb.client.Publish(ctx, channelName, data).Err(); err != nil {
		log.Printf("❌ Failed to publish event [%s]: %v", e.Type, err)
		return err
	}

	log.Printf("📢 Published event: %s", e.Type)
	return nil
}

// Subscribe returns a channel that receives events matching the given types.
// If no event types are specified, all events are received.
// Cancel the context to unsubscribe.
func (eb *EventBus) Subscribe(ctx context.Context, eventTypes ...string) (<-chan event.Event, error) {
	pubsub := eb.client.Subscribe(ctx, channelName)

	// Wait for confirmation that subscription is created
	if _, err := pubsub.Receive(ctx); err != nil {
		return nil, err
	}

	// Build a set of event types to filter on
	typeFilter := make(map[string]bool, len(eventTypes))
	for _, t := range eventTypes {
		typeFilter[t] = true
	}

	out := make(chan event.Event, 100)

	go func() {
		defer close(out)
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				log.Println("🔌 Subscriber disconnected")
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}

				var e event.Event
				if err := json.Unmarshal([]byte(msg.Payload), &e); err != nil {
					log.Printf("⚠️ Failed to unmarshal event: %v", err)
					continue
				}

				// Filter by event type if types were specified
				if len(typeFilter) > 0 && !typeFilter[e.Type] {
					continue
				}

				// Non-blocking send to prevent slow consumers from blocking others
				select {
				case out <- e:
				default:
					log.Printf("⚠️ Subscriber buffer full, dropping event: %s", e.Type)
				}
			}
		}
	}()

	log.Printf("👂 New subscriber (filter: %v)", eventTypes)
	return out, nil
}

// Close shuts down the event bus
func (eb *EventBus) Close() error {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.closed = true
	log.Println("Event bus closed")
	return nil
}
