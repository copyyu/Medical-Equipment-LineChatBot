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
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
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
	}
}
