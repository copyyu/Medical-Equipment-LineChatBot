package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupHealthRoutes configures health check endpoints
func SetupHealthRoutes(app *fiber.App) {
	// Direct health endpoint (basic check)
	app.Get("/health", handlers.HealthCheck)

	// API v1 health routes
	// v1 := app.Group("/api/v1")
	// health := v1.Group("/health")
	// {
	// 	// Basic health check
	// 	health.Get("/", handlers.HealthCheck)
	// }
}
