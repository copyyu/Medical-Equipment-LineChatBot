package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupNotificationRoutes configures notification endpoints
func SetupNotificationRoutes(app *fiber.App, notificationHandler *handlers.NotificationHandler) {
	// api := app.Group("/api/v1")
	// notifications := api.Group("/notifications")

	// Manual trigger endpoints
	app.Post("/send/june", notificationHandler.SendJuneAlerts)
	app.Post("/send/august", notificationHandler.SendAugustAlerts)

	// Summary
	app.Get("/summary", notificationHandler.GetSummary)

	// Settings
	app.Put("/settings", notificationHandler.UpdateSettings)
}
