package handlers

import (
	"medical-webhook/internal/infrastructure/database"
	redisinfra "medical-webhook/internal/infrastructure/redis"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck returns basic server health status (legacy endpoint, kept for compatibility)
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "healthy",
		"service": "medical-equipment-webhook",
		"version": "1.0.0",
	})
}

// LivenessCheck returns whether the application process is alive.
// Used by orchestrators (Kubernetes, Docker) to determine if the container
// needs to be restarted. This should NEVER check external dependencies.
func LivenessCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "alive",
	})
}

// ReadinessCheck returns whether the application is ready to accept traffic.
// Checks all critical dependencies (database, Redis).
// Used by load balancers to decide whether to route traffic to this instance.
func ReadinessCheck(c *fiber.Ctx) error {
	checks := fiber.Map{}
	ready := true

	// Check database
	if err := database.HealthCheck(); err != nil {
		checks["database"] = fiber.Map{"status": "down", "error": err.Error()}
		ready = false
	} else {
		checks["database"] = fiber.Map{"status": "up"}
	}

	// Check Redis (optional — degraded but not unready if Redis is down)
	redisClient := redisinfra.GetClient()
	if redisClient != nil {
		if err := redisClient.Ping(c.Context()).Err(); err != nil {
			checks["redis"] = fiber.Map{"status": "down", "error": err.Error()}
			// Redis is optional — don't mark unready, just degraded
		} else {
			checks["redis"] = fiber.Map{"status": "up"}
		}
	} else {
		checks["redis"] = fiber.Map{"status": "not_configured"}
	}

	status := fiber.StatusOK
	statusText := "ready"
	if !ready {
		status = fiber.StatusServiceUnavailable
		statusText = "not_ready"
	}

	return c.Status(status).JSON(fiber.Map{
		"status":  statusText,
		"service": "medical-equipment-webhook",
		"checks":  checks,
	})
}
