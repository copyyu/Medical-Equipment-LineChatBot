// Package logger configures the application's structured logger (log/slog).
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// parseLevel maps a textual level to slog.Level, defaulting to Info.
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

// isProduction reports whether env denotes a production-like environment.
func isProduction(env string) bool {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "production", "prod":
		return true
	default:
		return false
	}
}

// Init builds a structured slog.Logger for the given environment and level and
// installs it as the default. In production it emits JSON (machine-parseable);
// otherwise it emits human-friendly text. Installing it as the default also
// bridges the standard library's log package, so existing log.Printf calls are
// emitted through the same handler.
func Init(env, level string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(level)}

	var handler slog.Handler
	if isProduction(env) {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}
