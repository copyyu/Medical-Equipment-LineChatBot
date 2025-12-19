package main

import (
	"log"
	"medical-webhook/config"
	"medical-webhook/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create webhook handler
	webhookHandler, err := handlers.NewWebhookHandler(cfg.LineChannelToken, cfg.LineChannelSecret)
	if err != nil {
		log.Fatalf("❌ Failed to create webhook handler: %v", err)
	}

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Medical Equipment Webhook",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🏥 Medical Equipment Webhook Server",
			"status":  "running",
		})
	})

	app.Get("/health", handlers.HealthCheck)
	app.Post("/webhook", webhookHandler.HandleCallback)
	app.Post("/callback", webhookHandler.HandleCallback) // Alias

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
