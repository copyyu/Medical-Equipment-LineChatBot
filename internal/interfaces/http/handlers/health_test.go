package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestLivenessCheck_OK(t *testing.T) {
	app := fiber.New()
	app.Get("/livez", LivenessCheck)

	resp, err := app.Test(httptest.NewRequest("GET", "/livez", nil))
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("liveness status = %d, want 200", resp.StatusCode)
	}
}

func TestReadinessCheck_DBDownReturns503(t *testing.T) {
	// No database is connected in the test process, so readiness must report
	// not-ready with 503 (Redis is optional and reported as "disabled").
	app := fiber.New()
	app.Get("/readyz", ReadinessCheck)

	resp, err := app.Test(httptest.NewRequest("GET", "/readyz", nil))
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusServiceUnavailable {
		t.Fatalf("readiness status = %d, want 503 when DB is down", resp.StatusCode)
	}
}
