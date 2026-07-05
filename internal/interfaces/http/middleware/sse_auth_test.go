package middleware

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"medical-webhook/internal/application/dto"

	"github.com/gofiber/fiber/v2"
)

// stubAdminUsecase validates only the token "good".
type stubAdminUsecase struct{}

func (stubAdminUsecase) Register(context.Context, *dto.RegisterRequest) (*dto.AdminDetail, error) {
	return nil, nil
}
func (stubAdminUsecase) Login(context.Context, *dto.LoginRequest, string) (*dto.LoginResponse, error) {
	return nil, nil
}
func (stubAdminUsecase) Logout(context.Context, string) error { return nil }
func (stubAdminUsecase) GetProfile(context.Context, string) (*dto.AdminDetail, error) {
	return nil, nil
}
func (stubAdminUsecase) UpdateProfile(context.Context, string, *dto.UpdateProfileRequest) error {
	return nil
}
func (stubAdminUsecase) ChangePassword(context.Context, string, *dto.ChangePasswordRequest) error {
	return nil
}
func (stubAdminUsecase) ValidateToken(_ context.Context, token string) (*dto.AdminDetail, error) {
	if token == "good" {
		return &dto.AdminDetail{ID: "1", Username: "u", Role: "admin"}, nil
	}
	return nil, errors.New("invalid token")
}

func sseStatus(t *testing.T, target string, authHeader string) int {
	t.Helper()
	app := fiber.New()
	app.Get("/api/events/stream", SSEAuthMiddleware(stubAdminUsecase{}), func(c *fiber.Ctx) error {
		return c.SendString("stream")
	})
	req := httptest.NewRequest("GET", target, nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	return resp.StatusCode
}

func TestSSEAuthMiddleware(t *testing.T) {
	cases := []struct {
		name   string
		target string
		header string
		want   int
	}{
		{"valid token via query", "/api/events/stream?token=good", "", fiber.StatusOK},
		{"valid token via bearer header", "/api/events/stream", "Bearer good", fiber.StatusOK},
		{"no token", "/api/events/stream", "", fiber.StatusUnauthorized},
		{"invalid token via query", "/api/events/stream?token=bad", "", fiber.StatusUnauthorized},
		{"malformed header", "/api/events/stream", "Token good", fiber.StatusUnauthorized},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := sseStatus(t, tc.target, tc.header); got != tc.want {
				t.Errorf("status = %d, want %d", got, tc.want)
			}
		})
	}
}
