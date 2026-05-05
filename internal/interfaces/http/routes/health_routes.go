package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupHealthRoutes configures health check endpoints
func SetupHealthRoutes(app *fiber.App) {
	// Legacy health endpoint (kept for backward compatibility)
	app.Get("/health", handlers.HealthCheck)

	// Production health endpoints
	health := app.Group("/health")
	{
		// Liveness: is the process alive? (for container restart decisions)
		health.Get("/live", handlers.LivenessCheck)

		// Readiness: are dependencies ready? (for load balancer routing)
		health.Get("/ready", handlers.ReadinessCheck)
	}
}
