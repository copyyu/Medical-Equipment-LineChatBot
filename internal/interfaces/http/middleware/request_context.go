package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDKey is the Fiber Locals key holding the current request ID.
const RequestIDKey = "request_id"

type ctxKeyRequestID struct{}

// RequestContext ensures every request has a stable ID: it reuses an inbound
// X-Request-ID header when present, otherwise generates a UUID. The ID is echoed
// back in the response header and made available via Fiber Locals and the
// request's context (for downstream, context-aware logging).
func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get(fiber.HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set(fiber.HeaderXRequestID, reqID)
		c.Locals(RequestIDKey, reqID)
		c.SetUserContext(context.WithValue(c.UserContext(), ctxKeyRequestID{}, reqID))
		return c.Next()
	}
}

// RequestIDFromContext returns the request ID stored in ctx, or "" if absent.
func RequestIDFromContext(ctx context.Context) string {
	if ctx != nil {
		if id, ok := ctx.Value(ctxKeyRequestID{}).(string); ok {
			return id
		}
	}
	return ""
}
