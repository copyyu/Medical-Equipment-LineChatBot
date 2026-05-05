package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	AppEnv            string // dev, staging, prod
	LogLevel          string // debug, info, warn, error
	LineChannelToken  string
	LineChannelSecret string
	Port              string
	OCRURL            string
	RedisURL          string
	BaseURL           string
	DB                DatabaseConfig
	Contact           ContactConfig
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// ContactConfig holds contact information (loaded from env to avoid hardcoding)
type ContactConfig struct {
	CenterName     string
	Phone          string
	Email          string
	EmergencyPhone string
	WorkingHours   string
}

// Load returns configuration from environment variables or defaults
func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	return &Config{
		AppEnv:            getEnvOrDefault("APP_ENV", "dev"),
		LogLevel:          getEnvOrDefault("LOG_LEVEL", "info"),
		LineChannelToken:  os.Getenv("LINE_CHANNEL_TOKEN"),
		LineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		Port:              getEnvOrDefault("PORT", "3000"),
		OCRURL:            os.Getenv("OCR_API_URL"),
		RedisURL:          getEnvOrDefault("REDIS_URL", "redis://localhost:6379"),
		BaseURL:           os.Getenv("BASE_URL"),
		DB: DatabaseConfig{
			Host:            os.Getenv("DB_HOST"),
			Port:            os.Getenv("DB_PORT"),
			User:            os.Getenv("DB_USER"),
			Password:        os.Getenv("DB_PASSWORD"),
			Name:            os.Getenv("DB_NAME"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 30)) * time.Minute,
		},
		Contact: ContactConfig{
			CenterName:     getEnvOrDefault("CONTACT_CENTER_NAME", "ศูนย์เครื่องมือแพทย์"),
			Phone:          getEnvOrDefault("CONTACT_PHONE", ""),
			Email:          getEnvOrDefault("CONTACT_EMAIL", ""),
			EmergencyPhone: getEnvOrDefault("CONTACT_EMERGENCY_PHONE", ""),
			WorkingHours:   getEnvOrDefault("CONTACT_WORKING_HOURS", "จ-ศ 08:00-17:00"),
		},
	}
}

// Validate checks that all required configuration values are present.
// Returns an error listing all missing variables so the operator can fix them in one pass.
func (c *Config) Validate() error {
	var missing []string

	// LINE API credentials (required for core functionality)
	if c.LineChannelToken == "" {
		missing = append(missing, "LINE_CHANNEL_TOKEN")
	}
	if c.LineChannelSecret == "" {
		missing = append(missing, "LINE_CHANNEL_SECRET")
	}

	// Database (required — app cannot function without it)
	if c.DB.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.DB.Port == "" {
		missing = append(missing, "DB_PORT")
	}
	if c.DB.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.DB.Name == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: [%s]", strings.Join(missing, ", "))
	}
	return nil
}

// IsDev returns true if the application is running in development mode.
func (c *Config) IsDev() bool {
	return strings.EqualFold(c.AppEnv, "dev")
}

// IsProd returns true if the application is running in production mode.
func (c *Config) IsProd() bool {
	return strings.EqualFold(c.AppEnv, "prod") || strings.EqualFold(c.AppEnv, "production")
}

// getEnvOrDefault returns the environment variable value or a default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt returns the environment variable as int or a default
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
