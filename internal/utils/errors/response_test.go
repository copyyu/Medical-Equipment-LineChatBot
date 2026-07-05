package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestMapErrorToResponse_Sentinels(t *testing.T) {
	cases := []struct {
		err        error
		wantStatus int
		wantCode   string
	}{
		{ErrNotFound, fiber.StatusNotFound, "NOT_FOUND"},
		{ErrUnauthorized, fiber.StatusUnauthorized, "UNAUTHORIZED"},
		{ErrConflict, fiber.StatusConflict, "CONFLICT"},
		{ErrValidationFailed, fiber.StatusBadRequest, "VALIDATION_FAILED"},
		{ErrInvalidCredentials, fiber.StatusUnauthorized, "INVALID_CREDENTIALS"},
		{errors.New("something random"), fiber.StatusInternalServerError, "INTERNAL_ERROR"},
	}
	for _, tc := range cases {
		status, code, _ := MapErrorToResponse(tc.err)
		if status != tc.wantStatus || code != tc.wantCode {
			t.Errorf("MapErrorToResponse(%v) = (%d,%q), want (%d,%q)", tc.err, status, code, tc.wantStatus, tc.wantCode)
		}
	}
}

func TestMapErrorToResponse_RespectsFiberError(t *testing.T) {
	// A *fiber.Error should keep its own status instead of collapsing to 500.
	status, code, msg := MapErrorToResponse(fiber.NewError(fiber.StatusBadRequest, "Invalid request body"))
	if status != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want 400", status)
	}
	if code != "BAD_REQUEST" {
		t.Fatalf("code = %q, want BAD_REQUEST", code)
	}
	if msg != "Invalid request body" {
		t.Fatalf("msg = %q, want the fiber error message", msg)
	}
}

func TestMapErrorToResponse_WrappedSentinel(t *testing.T) {
	wrapped := fmt.Errorf("lookup failed: %w", ErrNotFound)
	status, code, _ := MapErrorToResponse(wrapped)
	if status != fiber.StatusNotFound || code != "NOT_FOUND" {
		t.Fatalf("wrapped ErrNotFound = (%d,%q), want (404,NOT_FOUND)", status, code)
	}
}
