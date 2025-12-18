package main

import (
	"log"
	"medical-webhook/config"
	"medical-webhook/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create webhook handler
	webhookHandler, err := handlers.NewWebhookHandler(cfg.LineChannelToken, cfg.LineChannelSecret)
	if err != nil {
		log.Fatalf("❌ Failed to create webhook handler: %v", err)
	}

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "🏥 Medical Equipment Webhook Server",
			"status":  "running",
		})
	})
	r.GET("/health", handlers.HealthCheck)
	r.POST("/webhook", webhookHandler.HandleCallback)
	r.POST("/callback", webhookHandler.HandleCallback) // Alias for LINE webhook

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
