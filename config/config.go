package config

import (
	"os"
)

// Config holds all configuration values
type Config struct {
	LineChannelToken  string
	LineChannelSecret string
	Port              string
}

// Load returns configuration from environment variables or defaults
func Load() *Config {
	return &Config{
		LineChannelToken:  getEnv("LINE_CHANNEL_TOKEN", "G2+t5K7oCV6uxnGMCbV8QoCBFUJ68MIqpm25bdc99oyhSVPiDFD/tPUPzqHCxKNqsXStBicQs0R1KQflpHbP/7Q+QKyfRhKflgc/UZq+bOIngj6rP3SGpxChpURom7THrEWYG+NuepetcIBKK5GMcAdB04t89/1O/w1cDnyilFU="),
		LineChannelSecret: getEnv("LINE_CHANNEL_SECRET", "4e28ea2ad3e766b64453193dd99aae49"),
		Port:              getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
