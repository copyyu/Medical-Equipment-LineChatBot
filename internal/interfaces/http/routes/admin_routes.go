package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App, adminHandler *handlers.AdminHandler, adminUsecase usecase.AdminUsecase) {
	// API routes
	api := app.Group("/api")

	// Admin routes
	admin := api.Group("/admin")

	// Public routes - ไม่ต้อง auth (rate-limited to blunt brute-force/abuse)
	authLimiter := middleware.AuthRateLimiter()
	admin.Post("/login", authLimiter, adminHandler.Login)
	admin.Post("/register", authLimiter, adminHandler.Register)

	// Protected routes - ต้อง auth (Bearer token)
	adminProtected := admin.Group("", middleware.AuthMiddleware(adminUsecase))
	adminProtected.Post("/logout", adminHandler.Logout)
	// adminProtected.Get("/profile", adminHandler.GetProfile)
	// adminProtected.Put("/profile", adminHandler.UpdateProfile)
	// adminProtected.Post("/change-password", adminHandler.ChangePassword)
}
