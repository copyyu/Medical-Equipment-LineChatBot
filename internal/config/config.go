package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	LineChannelToken  string
	LineChannelSecret string
	Port              string
	OCRURL            string
	DB                DatabaseConfig
	Contact           ContactConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
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
		LineChannelToken:  os.Getenv("LINE_CHANNEL_TOKEN"),
		LineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		Port:              os.Getenv("PORT"),
		OCRURL:            os.Getenv("OCR_API_URL"),
		DB: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
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

// getEnvOrDefault returns the environment variable value or a default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
