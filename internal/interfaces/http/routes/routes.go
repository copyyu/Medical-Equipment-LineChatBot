package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// Setup configures all application routes
func Setup(app *fiber.App, webhookHandler *handlers.WebhookHandler) {
	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "🏥 Medical Equipment Webhook Server",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// Health check
	app.Get("/health", handlers.HealthCheck)

	// Webhook endpoints
	app.Post("/webhook", webhookHandler.HandleCallback)
	app.Post("/callback", webhookHandler.HandleCallback)
}
