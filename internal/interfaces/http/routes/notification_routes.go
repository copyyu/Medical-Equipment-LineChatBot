package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupNotificationRoutes - protected routes for notifications
func SetupNotificationRoutes(app *fiber.App, notificationHandler *handlers.NotificationHandler, adminUsecase usecase.AdminUsecase) {

	// The Excel export link is embedded as an "Export" button in the LINE Flex
	// alerts and is opened in the user's browser without a bearer token, so it
	// cannot sit behind AuthMiddleware. Access is instead gated by an HMAC
	// signature + expiry on the URL, validated inside the handler (see exporturl).
	app.Get("/notifications/export/expiry", notificationHandler.DownloadExpiryExcel)

	// Notification routes - protected
	notifGroup := app.Group("", middleware.AuthMiddleware(adminUsecase))

	// Manual trigger endpoints
	notifGroup.Post("/send/june", notificationHandler.SendJuneAlerts)
	notifGroup.Post("/send/august", notificationHandler.SendAugustAlerts)

	// Manual trigger endpoints for testing cronjob broadcasts. These send real
	// LINE messages to real recipients, so they must require authentication.
	notifGroup.Get("/test/cron/june", notificationHandler.TestJuneAlerts)
	notifGroup.Get("/test/cron/august", notificationHandler.TestAugustAlerts)

	// Summary
	notifGroup.Get("/summary", notificationHandler.GetSummary)

	// Settings
	notifGroup.Put("/settings", notificationHandler.UpdateSettings)
}
