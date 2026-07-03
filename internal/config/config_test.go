package config

import (
	"strings"
	"testing"
)

func validConfig() *Config {
	c := &Config{
		LineChannelToken:  "token",
		LineChannelSecret: "secret",
	}
	c.DB.Host = "localhost"
	c.DB.Port = "5432"
	c.DB.User = "postgres"
	c.DB.Name = "medical"
	return c
}

func TestValidate_OK(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("expected valid config to pass, got: %v", err)
	}
}

func TestValidate_MissingSecretReported(t *testing.T) {
	c := validConfig()
	c.LineChannelSecret = ""
	err := c.Validate()
	if err == nil {
		t.Fatalf("expected error for empty LINE_CHANNEL_SECRET")
	}
	if !strings.Contains(err.Error(), "LINE_CHANNEL_SECRET") {
		t.Fatalf("error should name the missing var, got: %v", err)
	}
}

func TestValidate_WhitespaceCountsAsMissing(t *testing.T) {
	c := validConfig()
	c.DB.Host = "   "
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "DB_HOST") {
		t.Fatalf("whitespace-only value should be treated as missing, got: %v", err)
	}
}

func TestValidate_ReportsAllMissing(t *testing.T) {
	c := &Config{}
	err := c.Validate()
	if err == nil {
		t.Fatalf("expected error for empty config")
	}
	for _, name := range []string{"LINE_CHANNEL_TOKEN", "LINE_CHANNEL_SECRET", "DB_HOST", "DB_PORT", "DB_USER", "DB_NAME"} {
		if !strings.Contains(err.Error(), name) {
			t.Fatalf("expected %s to be reported missing, got: %v", name, err)
		}
	}
}
