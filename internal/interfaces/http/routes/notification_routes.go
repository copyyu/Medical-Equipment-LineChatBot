package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupNotificationRoutes registers notification endpoints. All of them require
// admin authentication — they either send LINE notifications or expose internal
// data, so none may be public.
func SetupNotificationRoutes(app *fiber.App, notificationHandler *handlers.NotificationHandler, adminUsecase usecase.AdminUsecase) {
	notifGroup := app.Group("", middleware.AuthMiddleware(adminUsecase))

	// Export (exposes equipment/expiry data)
	notifGroup.Get("/notifications/export/expiry", notificationHandler.DownloadExpiryExcel)

	// Manual trigger endpoints (send real LINE notifications)
	notifGroup.Post("/send/june", notificationHandler.SendJuneAlerts)
	notifGroup.Post("/send/august", notificationHandler.SendAugustAlerts)

	// Manual test triggers (kept for diagnostics, now authenticated)
	notifGroup.Get("/test/cron/june", notificationHandler.TestJuneAlerts)
	notifGroup.Get("/test/cron/august", notificationHandler.TestAugustAlerts)

	// Summary
	notifGroup.Get("/summary", notificationHandler.GetSummary)

	// Settings
	notifGroup.Put("/settings", notificationHandler.UpdateSettings)
}
