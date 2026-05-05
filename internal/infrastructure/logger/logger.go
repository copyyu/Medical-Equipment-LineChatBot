// Package logger provides a centralized, structured logging facility using Go's
// log/slog package. It supports JSON output for production and text output for
// development, configurable via the LOG_LEVEL and APP_ENV environment variables.
package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// contextKey is an unexported type for context keys in this package.
type contextKey string

const (
	// RequestIDKey is the context key for the request ID.
	RequestIDKey contextKey = "request_id"
	// AdminIDKey is the context key for the admin ID.
	AdminIDKey contextKey = "admin_id"
	// AdminUsernameKey is the context key for the admin username.
	AdminUsernameKey contextKey = "admin_username"
)

// Setup initializes the global slog logger.
// Call this once at application startup.
//
//   - appEnv: "dev" uses human-readable text output; anything else uses JSON.
//   - logLevel: one of "debug", "info", "warn", "error" (default: "info").
func Setup(appEnv, logLevel string) {
	level := parseLevel(logLevel)

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug, // include source file info in debug mode
	}

	if strings.EqualFold(appEnv, "dev") {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

// FromContext returns a logger enriched with request-scoped values
// (request_id, admin_id, etc.) stored in the context.
func FromContext(ctx context.Context) *slog.Logger {
	logger := slog.Default()

	if rid, ok := ctx.Value(RequestIDKey).(string); ok && rid != "" {
		logger = logger.With("request_id", rid)
	}
	if adminID, ok := ctx.Value(AdminIDKey).(string); ok && adminID != "" {
		logger = logger.With("admin_id", adminID)
	}
	if username, ok := ctx.Value(AdminUsernameKey).(string); ok && username != "" {
		logger = logger.With("admin_username", username)
	}

	return logger
}

// WithRequestID returns a new context with the given request ID attached.
func WithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, RequestIDKey, rid)
}

// WithAdmin returns a new context with admin information attached.
func WithAdmin(ctx context.Context, adminID, username string) context.Context {
	ctx = context.WithValue(ctx, AdminIDKey, adminID)
	ctx = context.WithValue(ctx, AdminUsernameKey, username)
	return ctx
}

// parseLevel converts a string level to slog.Level.
func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
