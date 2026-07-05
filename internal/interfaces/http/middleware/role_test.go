package middleware

import (
	"io"
	"net/http/httptest"
	"testing"

	"medical-webhook/internal/domain/line/entity"

	"github.com/gofiber/fiber/v2"
)

// newRoleApp builds an app whose /x route injects the given role (as
// AuthMiddleware would) and then guards with RequireRole(allowed...).
func newRoleApp(injectRole string, hasRole bool, allowed ...entity.AdminRole) *fiber.App {
	app := fiber.New()
	app.Get("/x",
		func(c *fiber.Ctx) error {
			if hasRole {
				c.Locals("admin_role", injectRole)
			}
			return c.Next()
		},
		RequireRole(allowed...),
		func(c *fiber.Ctx) error { return c.SendString("ok") },
	)
	return app
}

func status(t *testing.T, app *fiber.App) int {
	t.Helper()
	resp, err := app.Test(httptest.NewRequest("GET", "/x", nil))
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	return resp.StatusCode
}

func TestRequireRole(t *testing.T) {
	cases := []struct {
		name    string
		role    string
		hasRole bool
		allowed []entity.AdminRole
		want    int
	}{
		{"super admin allowed", string(entity.RoleSuperAdmin), true, []entity.AdminRole{entity.RoleSuperAdmin}, fiber.StatusOK},
		{"regular admin forbidden", string(entity.RoleAdmin), true, []entity.AdminRole{entity.RoleSuperAdmin}, fiber.StatusForbidden},
		{"staff forbidden", string(entity.RoleStaff), true, []entity.AdminRole{entity.RoleSuperAdmin}, fiber.StatusForbidden},
		{"admin allowed when admin permitted", string(entity.RoleAdmin), true, []entity.AdminRole{entity.RoleSuperAdmin, entity.RoleAdmin}, fiber.StatusOK},
		{"missing role unauthorized", "", false, []entity.AdminRole{entity.RoleSuperAdmin}, fiber.StatusUnauthorized},
		{"empty role unauthorized", "", true, []entity.AdminRole{entity.RoleSuperAdmin}, fiber.StatusUnauthorized},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newRoleApp(tc.role, tc.hasRole, tc.allowed...)
			if got := status(t, app); got != tc.want {
				t.Errorf("status = %d, want %d", got, tc.want)
			}
		})
	}
}
