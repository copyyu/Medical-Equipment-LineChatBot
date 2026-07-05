package middleware

import (
	"medical-webhook/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
)

// RequireRole authorizes requests whose authenticated admin holds one of the
// allowed roles. It MUST run after AuthMiddleware, which populates the
// "admin_role" local. Responds 403 (FORBIDDEN) when the role is not permitted.
//
// Example:
//
//	adminProtected.Delete("/admins/:id",
//	    middleware.RequireRole(string(entity.RoleSuperAdmin)),
//	    adminHandler.DeleteAdmin)
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("admin_role").(string)
		if _, ok := allowed[role]; !ok {
			return errors.Error(c, fiber.NewError(fiber.StatusForbidden, "Insufficient permissions"))
		}
		return c.Next()
	}
}
