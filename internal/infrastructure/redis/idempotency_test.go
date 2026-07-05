package redis

import (
	"context"
	"testing"
	"time"
)

func TestNewIdempotencyStore_DefaultTTL(t *testing.T) {
	s := NewIdempotencyStore(nil, 0)
	if s.ttl != defaultTTL {
		t.Fatalf("expected default TTL %v, got %v", defaultTTL, s.ttl)
	}
}

func TestNewIdempotencyStore_CustomTTL(t *testing.T) {
	s := NewIdempotencyStore(nil, 5*time.Minute)
	if s.ttl != 5*time.Minute {
		t.Fatalf("expected 5m TTL, got %v", s.ttl)
	}
}

// MarkProcessed must fail open (return true = "process it") whenever the store
// cannot reach Redis, so events are never silently dropped.

func TestMarkProcessed_FailsOpenWithoutClient(t *testing.T) {
	s := NewIdempotencyStore(nil, 0)
	ok, err := s.MarkProcessed(context.Background(), "evt-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected fail-open (true) when the client is unavailable")
	}
}

func TestMarkProcessed_NilReceiverFailsOpen(t *testing.T) {
	var s *IdempotencyStore
	ok, err := s.MarkProcessed(context.Background(), "evt-1")
	if err != nil || !ok {
		t.Fatalf("nil store should fail open, got ok=%v err=%v", ok, err)
	}
}

func TestMarkProcessed_EmptyEventIDFailsOpen(t *testing.T) {
	s := NewIdempotencyStore(nil, 0)
	ok, err := s.MarkProcessed(context.Background(), "")
	if err != nil || !ok {
		t.Fatalf("empty event id should fail open, got ok=%v err=%v", ok, err)
	}
}
