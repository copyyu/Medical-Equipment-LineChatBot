package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupWebhookRoutes configures webhook endpoints
func SetupWebhookRoutes(app *fiber.App, webhookHandler *handlers.WebhookHandler) {
	// Direct webhook endpoints (for LINE Platform configuration)
	app.Post("/webhook", webhookHandler.HandleCallback)
	app.Post("/callback", webhookHandler.HandleCallback)

	// API v1 webhook routes
	// v1 := app.Group("/api/v1")
	// webhook := v1.Group("/webhook")
	// {
	// 	webhook.Post("/", webhookHandler.HandleCallback)
	// 	webhook.Post("/callback", webhookHandler.HandleCallback) // Alias
	// }
}
