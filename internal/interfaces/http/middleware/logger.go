package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// StructuredLogger is a Fiber middleware that logs every request using slog
// with structured fields: method, path, status, latency, ip, request_id, etc.
func StructuredLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		// Build structured log fields
		attrs := []slog.Attr{
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.String("latency", latency.String()),
			slog.String("ip", c.IP()),
			slog.String("user_agent", c.Get("User-Agent")),
		}

		// Add request ID if available
		if rid, ok := c.Locals("request_id").(string); ok && rid != "" {
			attrs = append(attrs, slog.String("request_id", rid))
		}

		// Add admin info if authenticated
		if adminID, ok := c.Locals("admin_id").(interface{ String() string }); ok {
			attrs = append(attrs, slog.String("admin_id", adminID.String()))
		}
		if username, ok := c.Locals("admin_username").(string); ok && username != "" {
			attrs = append(attrs, slog.String("admin_username", username))
		}

		// Add error if present
		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
		}

		// Choose log level based on status code
		args := make([]any, len(attrs))
		for i, a := range attrs {
			args[i] = a
		}

		switch {
		case status >= 500:
			slog.Error("HTTP Request", args...)
		case status >= 400:
			slog.Warn("HTTP Request", args...)
		default:
			slog.Info("HTTP Request", args...)
		}

		return err
	}
}
