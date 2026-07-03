package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequestLogger emits one structured access-log line per request via slog,
// including the request ID, method, path, status, latency and client IP. The
// log level scales with the response status (5xx=error, 4xx=warn, else info).
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		// Derive the status: when a handler returns an error the app ErrorHandler
		// runs after this middleware, so read the intended code from the error.
		status := c.Response().StatusCode()
		if err != nil {
			if fe, ok := err.(*fiber.Error); ok {
				status = fe.Code
			} else {
				status = fiber.StatusInternalServerError
			}
		}

		level := slog.LevelInfo
		switch {
		case status >= 500:
			level = slog.LevelError
		case status >= 400:
			level = slog.LevelWarn
		}

		reqID, _ := c.Locals(RequestIDKey).(string)
		attrs := []slog.Attr{
			slog.String("request_id", reqID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.Int64("latency_ms", time.Since(start).Milliseconds()),
			slog.String("ip", c.IP()),
		}
		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
		}
		slog.LogAttrs(c.UserContext(), level, "http_request", attrs...)

		return err
	}
}
