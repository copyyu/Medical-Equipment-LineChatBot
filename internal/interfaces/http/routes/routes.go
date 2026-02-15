package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// Setup configures all application routes
func Setup(app *fiber.App, webhookHandler *handlers.WebhookHandler,
	notificationHandler *handlers.NotificationHandler,
	equipmentImportHandler *handlers.EquipmentImportHandler,
	adminHandler *handlers.AdminHandler,
	dashboardHandler *handlers.DashboardHandler,
	equipmentHandler *handlers.EquipmentHandler,
	ticketHandler *handlers.TicketHandler,
	adminUsecase usecase.AdminUsecase) {
	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "🏥 Medical Equipment Webhook Server",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// Setup health routes (public)
	SetupHealthRoutes(app)

	// Setup webhook routes (public - verified by LINE channel secret)
	SetupWebhookRoutes(app, webhookHandler)

	// Setup admin routes (login = public, rest = protected)
	SetupAdminRoutes(app, adminHandler, adminUsecase)

	// ===== Protected routes (require auth) =====

	// Setup notifications routes
	SetupNotificationRoutes(app, notificationHandler, adminUsecase)

	// Setup equipmentImport routes
	SetupEquipmentImportRoutes(app, equipmentImportHandler, adminUsecase)

	// Setup dashboard routes
	SetupDashboardRoutes(app, dashboardHandler, adminUsecase)

	// Setup equipment routes
	SetupEquipmentRoutes(app, equipmentHandler, adminUsecase)

	// Setup ticket routes
	SetupTicketRoutes(app, ticketHandler, adminUsecase)

	// 404 handler (must be last)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Route not found",
			"success": false,
			"path":    c.Path(),
		})
	})
}
