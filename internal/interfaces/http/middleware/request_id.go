package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

// RequestID generates a unique request ID for every incoming request.
// The ID is stored in c.Locals("request_id") and set as a response header.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Use client-supplied ID if present (for distributed tracing),
		// otherwise generate a new one.
		rid := c.Get(RequestIDHeader)
		if rid == "" {
			rid = uuid.New().String()
		}

		// Store in Locals for downstream handlers
		c.Locals("request_id", rid)

		// Set response header so callers can correlate
		c.Set(RequestIDHeader, rid)

		return c.Next()
	}
}
