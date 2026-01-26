package middleware

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/utils/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware - ใช้สำหรับ routes ที่ต้อง authentication
func AuthMiddleware(adminUsecase usecase.AdminUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.Unauthorized(c, "No authorization token provided")
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errors.Unauthorized(c, "Invalid authorization format. Use: Bearer <token>")
		}

		token := parts[1]
		if token == "" {
			return errors.Unauthorized(c, "Token is empty")
		}

		// Validate token
		admin, err := adminUsecase.ValidateToken(c.Context(), token)
		if err != nil {
			return errors.Unauthorized(c, "Invalid or expired token")
		}

		// Store admin info in context for use in handlers
		c.Locals("admin_id", admin.ID)
		c.Locals("admin_username", admin.Username)
		c.Locals("admin_email", admin.Email)
		c.Locals("admin", admin)
		c.Locals("admin_role", admin.Role)

		// Continue to next handler
		return c.Next()
	}
}
