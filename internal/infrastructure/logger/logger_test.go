package logger

import (
	"log/slog"
	"testing"
)

func TestParseLevel(t *testing.T) {
	cases := map[string]slog.Level{
		"debug":    slog.LevelDebug,
		"DEBUG":    slog.LevelDebug,
		"info":     slog.LevelInfo,
		"":         slog.LevelInfo,
		"nonsense": slog.LevelInfo,
		"warn":     slog.LevelWarn,
		"warning":  slog.LevelWarn,
		"error":    slog.LevelError,
		" Error ":  slog.LevelError,
	}
	for in, want := range cases {
		if got := parseLevel(in); got != want {
			t.Errorf("parseLevel(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestIsProduction(t *testing.T) {
	for _, s := range []string{"production", "PROD", "Prod"} {
		if !isProduction(s) {
			t.Errorf("isProduction(%q) = false, want true", s)
		}
	}
	for _, s := range []string{"development", "dev", "", "staging"} {
		if isProduction(s) {
			t.Errorf("isProduction(%q) = true, want false", s)
		}
	}
}

func TestInit_ReturnsLoggerAndSetsDefault(t *testing.T) {
	l := Init("development", "debug")
	if l == nil {
		t.Fatal("Init returned nil logger")
	}
	if slog.Default() == nil {
		t.Fatal("Init did not set a default logger")
	}
}
