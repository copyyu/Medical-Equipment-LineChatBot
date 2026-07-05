package handlers

import (
	"medical-webhook/internal/infrastructure/database"
	redisinfra "medical-webhook/internal/infrastructure/redis"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck returns a basic server status (kept for backward compatibility).
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "healthy",
		"service": "medical-equipment-webhook",
		"version": "1.0.0",
	})
}

// LivenessCheck reports whether the process is up. It performs no dependency
// checks, so an orchestrator uses it only to decide whether to restart the pod.
func LivenessCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "alive"})
}

// ReadinessCheck reports whether the service can serve traffic. The database is
// required (its failure returns 503); Redis is reported but treated as optional
// because the app is designed to run in a degraded mode without it.
func ReadinessCheck(c *fiber.Ctx) error {
	checks := fiber.Map{}
	ready := true

	if err := database.HealthCheck(); err != nil {
		checks["database"] = "unavailable"
		ready = false
	} else {
		checks["database"] = "ok"
	}

	switch client := redisinfra.GetClient(); {
	case client == nil:
		checks["redis"] = "disabled"
	default:
		if err := client.Ping(c.UserContext()).Err(); err != nil {
			checks["redis"] = "unavailable" // informational; does not fail readiness
		} else {
			checks["redis"] = "ok"
		}
	}

	status := fiber.StatusOK
	readiness := "ready"
	if !ready {
		status = fiber.StatusServiceUnavailable
		readiness = "not_ready"
	}

	return c.Status(status).JSON(fiber.Map{
		"status": readiness,
		"checks": checks,
	})
}
