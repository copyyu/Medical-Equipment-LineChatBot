package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App, adminHandler *handlers.AdminHandler, adminUsecase usecase.AdminUsecase) {
	// API routes
	api := app.Group("/api")

	// Admin routes
	admin := api.Group("/admin")

	// Public routes - only login is public (rate-limited to blunt brute-force)
	admin.Post("/login", middleware.AuthRateLimiter(), adminHandler.Login)

	// Protected routes - ต้อง auth (Bearer token)
	adminProtected := admin.Group("", middleware.AuthMiddleware(adminUsecase))
	// Creating admins is restricted to super-admins. Bootstrap the first one via
	// the ADMIN_BOOTSTRAP_* env vars (see docs/CONFIGURATION.md).
	adminProtected.Post("/register", middleware.RequireRole(string(entity.RoleSuperAdmin)), adminHandler.Register)
	adminProtected.Post("/logout", adminHandler.Logout)

	// Self-service: an authenticated admin manages their own profile/password
	adminProtected.Get("/profile", adminHandler.GetProfile)
	adminProtected.Put("/profile", adminHandler.UpdateProfile)
	adminProtected.Post("/change-password", adminHandler.ChangePassword)
}
