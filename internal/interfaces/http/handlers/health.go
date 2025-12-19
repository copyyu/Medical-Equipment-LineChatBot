package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// HealthCheck returns server health status
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "healthy",
		"service": "medical-equipment-webhook",
		"version": "1.0.0",
	})
}
