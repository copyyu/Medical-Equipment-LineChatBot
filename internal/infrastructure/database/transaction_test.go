package database

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

// These tests cover the context <-> transaction plumbing that makes repositories
// transparently join a TxManager transaction. They intentionally use pointer
// identity (no real DB) to verify the routing logic in isolation.

func TestDBFromContext_NoTxReturnsFallback(t *testing.T) {
	fallback := &gorm.DB{}
	if got := DBFromContext(context.Background(), fallback); got != fallback {
		t.Fatalf("expected fallback handle when context carries no tx")
	}
}

func TestDBFromContext_NilContextReturnsFallback(t *testing.T) {
	fallback := &gorm.DB{}
	//nolint:staticcheck // deliberately passing a nil context to exercise the guard
	if got := DBFromContext(nil, fallback); got != fallback {
		t.Fatalf("expected fallback handle for a nil context")
	}
}

func TestContextWithTx_RoundTrips(t *testing.T) {
	tx := &gorm.DB{}
	fallback := &gorm.DB{}
	ctx := ContextWithTx(context.Background(), tx)
	if got := DBFromContext(ctx, fallback); got != tx {
		t.Fatalf("expected the tx stored in context to be returned, got the fallback")
	}
}

func TestContextWithTx_NilTxFallsBack(t *testing.T) {
	fallback := &gorm.DB{}
	ctx := ContextWithTx(context.Background(), nil)
	if got := DBFromContext(ctx, fallback); got != fallback {
		t.Fatalf("expected fallback when the stored tx is nil")
	}
}
