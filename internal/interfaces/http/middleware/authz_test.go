package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func newRoleApp(role string) *fiber.App {
	app := fiber.New()
	// Simulate AuthMiddleware having set the admin role.
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("admin_role", role)
		return c.Next()
	})
	app.Get("/admin-only", RequireRole("admin", "super_admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	return app
}

func TestRequireRole_Allowed(t *testing.T) {
	resp, err := newRoleApp("admin").Test(httptest.NewRequest("GET", "/admin-only", nil))
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("allowed role got %d, want 200", resp.StatusCode)
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	resp, err := newRoleApp("staff").Test(httptest.NewRequest("GET", "/admin-only", nil))
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("disallowed role got %d, want 403", resp.StatusCode)
	}
}

func TestRequireRole_MissingRole(t *testing.T) {
	app := fiber.New() // no admin_role local set
	app.Get("/x", RequireRole("admin"), func(c *fiber.Ctx) error { return c.SendStatus(200) })
	resp, err := app.Test(httptest.NewRequest("GET", "/x", nil))
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("missing role got %d, want 403", resp.StatusCode)
	}
}
