package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// Setup configures all application routes
func Setup(app *fiber.App, webhookHandler *handlers.WebhookHandler, notificationHandler *handlers.NotificationHandler, equipmentImportHandler *handlers.EquipmentImportHandler) {
	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "🏥 Medical Equipment Webhook Server",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// Setup health routes
	SetupHealthRoutes(app)

	// Setup webhook routes
	SetupWebhookRoutes(app, webhookHandler)

	// Setup notifications routes
	SetupNotificationRoutes(app, notificationHandler)

	// Setup equipmentImport routes
	SetupEquipmentImportRoutes(app, equipmentImportHandler)
	// 404 handler (must be last)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Route not found",
			"success": false,
			"path":    c.Path(),
		})
	})
}
