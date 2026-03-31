package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupNotificationRoutes - protected routes for notifications
func SetupNotificationRoutes(app *fiber.App, notificationHandler *handlers.NotificationHandler, adminUsecase usecase.AdminUsecase) {

	app.Get("/notifications/export/expiry", notificationHandler.DownloadExpiryExcel)

	// Temporary manual trigger endpoints for testing cronjob messages
	app.Get("/test/cron/june", notificationHandler.TestJuneAlerts)
	app.Get("/test/cron/august", notificationHandler.TestAugustAlerts)

	// Notification routes - protected
	notifGroup := app.Group("", middleware.AuthMiddleware(adminUsecase))

	// Manual trigger endpoints
	notifGroup.Post("/send/june", notificationHandler.SendJuneAlerts)
	notifGroup.Post("/send/august", notificationHandler.SendAugustAlerts)

	// Summary
	notifGroup.Get("/summary", notificationHandler.GetSummary)

	// Settings
	notifGroup.Put("/settings", notificationHandler.UpdateSettings)
}
