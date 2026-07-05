package middleware

import (
	"context"
	"testing"
)

func TestRequestIDFromContext_RoundTrip(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKeyRequestID{}, "abc-123")
	if got := RequestIDFromContext(ctx); got != "abc-123" {
		t.Fatalf("expected abc-123, got %q", got)
	}
}

func TestRequestIDFromContext_Absent(t *testing.T) {
	if got := RequestIDFromContext(context.Background()); got != "" {
		t.Fatalf("expected empty string when no request id, got %q", got)
	}
	//nolint:staticcheck // exercise the nil-context guard
	if got := RequestIDFromContext(nil); got != "" {
		t.Fatalf("expected empty string for nil context, got %q", got)
	}
}
