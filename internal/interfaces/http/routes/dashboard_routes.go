package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupDashboardRoutes(app *fiber.App, dashboardHandler *handlers.DashboardHandler) {
	// API routes
	api := app.Group("/api")

	// Dashboard routes
	dashboard := api.Group("/dashboard")

	// Public routes (can add auth middleware later if needed)
	dashboard.Get("/summary", dashboardHandler.GetSummary)
}
