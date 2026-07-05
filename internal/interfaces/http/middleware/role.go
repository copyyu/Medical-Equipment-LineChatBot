package middleware

import (
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
)

// RequireRole restricts a route to admins whose role is in the allowed set. It
// must run after AuthMiddleware, which populates the "admin_role" local.
func RequireRole(allowed ...entity.AdminRole) fiber.Handler {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, r := range allowed {
		allowedSet[string(r)] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("admin_role").(string)
		if !ok || role == "" {
			return errors.Unauthorized(c, "Missing role in authenticated context")
		}
		if _, allowed := allowedSet[role]; !allowed {
			return errors.Forbidden(c, "Insufficient permissions for this action")
		}
		return c.Next()
	}
}
