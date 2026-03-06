package redis

import (
	"context"
	"fmt"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

var client *goredis.Client

// Connect initializes the Redis client connection
func Connect(redisURL string) error {
	opts, err := goredis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	c := goredis.NewClient(opts)

	// Test connection
	ctx := context.Background()
	if err := c.Ping(ctx).Err(); err != nil {
		c.Close()
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Only set global client after successful connection
	client = c

	log.Printf("✅ Connected to Redis: %s", redisURL)
	return nil
}

// GetClient returns the Redis client instance
func GetClient() *goredis.Client {
	return client
}

// Close closes the Redis connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
