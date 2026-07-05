package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const (
	idempotencyKeyPrefix = "webhook:event:"
	defaultTTL           = 24 * time.Hour
)

// IdempotencyStore records processed webhook event IDs in Redis so that
// duplicate deliveries (LINE redelivers with at-least-once semantics) are acted
// on at most once. It relies on an atomic SET NX with a TTL — callers must use
// MarkProcessed's return value as the sole gate; there is deliberately no
// separate "check" method that would reintroduce a check-then-act race.
type IdempotencyStore struct {
	client *goredis.Client
	ttl    time.Duration
}

// NewIdempotencyStore creates a new idempotency store. If ttl is 0 it defaults
// to 24 hours (comfortably longer than LINE's redelivery window).
func NewIdempotencyStore(client *goredis.Client, ttl time.Duration) *IdempotencyStore {
	if ttl == 0 {
		ttl = defaultTTL
	}
	return &IdempotencyStore{
		client: client,
		ttl:    ttl,
	}
}

// MarkProcessed atomically records eventID as processed and reports whether this
// call was the first to do so:
//   - (true, nil):  first time — the caller should process the event. This is
//     also returned (fail-open) when no store/client is configured or eventID is
//     empty, so a Redis outage degrades to "process anyway" rather than silently
//     dropping events.
//   - (false, nil): the event was already processed — a duplicate to skip.
//   - (_, err):     a Redis error occurred; the caller decides how to react
//     (the webhook handler logs and processes anyway).
func (s *IdempotencyStore) MarkProcessed(ctx context.Context, eventID string) (bool, error) {
	if s == nil || s.client == nil || eventID == "" {
		return true, nil
	}
	key := idempotencyKeyPrefix + eventID
	// SET NX — atomic check-and-set: ok is true only if the key did not exist.
	ok, err := s.client.SetNX(ctx, key, "1", s.ttl).Result()
	if err != nil {
		return false, fmt.Errorf("idempotency: failed to mark event %q as processed: %w", eventID, err)
	}
	return ok, nil
}
