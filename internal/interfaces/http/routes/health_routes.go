package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupHealthRoutes configures health check endpoints.
func SetupHealthRoutes(app *fiber.App) {
	// Basic health endpoint (kept for backward compatibility)
	app.Get("/health", handlers.HealthCheck)

	// Kubernetes-style probes
	app.Get("/livez", handlers.LivenessCheck)   // process is up
	app.Get("/readyz", handlers.ReadinessCheck) // dependencies (DB) are reachable
}
