package middleware

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/utils/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SSEAuthMiddleware authenticates Server-Sent Events connections. Browsers'
// EventSource API cannot set an Authorization header, so the token is accepted
// from the "token" query parameter as well as a Bearer header (for non-browser
// clients). Frontend usage: new EventSource(`/api/events/stream?token=${bearer}`).
func SSEAuthMiddleware(adminUsecase usecase.AdminUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			if parts := strings.SplitN(c.Get("Authorization"), " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
		if token == "" {
			return errors.Unauthorized(c, "No token provided (use ?token= for EventSource)")
		}

		admin, err := adminUsecase.ValidateToken(c.Context(), token)
		if err != nil {
			return errors.Unauthorized(c, "Invalid or expired token")
		}

		c.Locals("admin_id", admin.ID)
		c.Locals("admin_username", admin.Username)
		c.Locals("admin_role", admin.Role)
		return c.Next()
	}
}
