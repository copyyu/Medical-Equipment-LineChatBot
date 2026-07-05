package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupNotificationRoutes registers notification endpoints. The Excel export
// link is embedded as an "Export" button in the LINE Flex alerts and opened in
// the user's browser without a bearer token, so it stays public but is gated by
// an HMAC signature + expiry validated inside the handler (see exporturl). The
// remaining endpoints send LINE notifications or expose internal data, so they
// require admin authentication.
func SetupNotificationRoutes(app *fiber.App, notificationHandler *handlers.NotificationHandler, adminUsecase usecase.AdminUsecase) {
	// Public but signature-gated. MUST be registered before the AuthMiddleware
	// group below, which otherwise catches all subsequent routes.
	app.Get("/notifications/export/expiry", notificationHandler.DownloadExpiryExcel)

	// Protected — require admin authentication.
	notifGroup := app.Group("", middleware.AuthMiddleware(adminUsecase))

	// Manual trigger endpoints (send real LINE notifications)
	notifGroup.Post("/send/june", notificationHandler.SendJuneAlerts)
	notifGroup.Post("/send/august", notificationHandler.SendAugustAlerts)

	// Manual test triggers (send real LINE messages to real recipients)
	notifGroup.Get("/test/cron/june", notificationHandler.TestJuneAlerts)
	notifGroup.Get("/test/cron/august", notificationHandler.TestAugustAlerts)

	// Summary
	notifGroup.Get("/summary", notificationHandler.GetSummary)

	// Settings
	notifGroup.Put("/settings", notificationHandler.UpdateSettings)
}
