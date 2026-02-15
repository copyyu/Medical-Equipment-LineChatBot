package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupDashboardRoutes(app *fiber.App, dashboardHandler *handlers.DashboardHandler, adminUsecase usecase.AdminUsecase) {
	api := app.Group("/api")

	// Dashboard routes - protected
	dashboard := api.Group("/dashboard", middleware.AuthMiddleware(adminUsecase))
	dashboard.Get("/summary", dashboardHandler.GetSummary)
}
