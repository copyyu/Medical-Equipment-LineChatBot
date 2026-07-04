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

	// Public routes - ไม่ต้อง auth
	admin.Post("/login", adminHandler.Login)

	// Protected routes - ต้อง auth (Bearer token)
	adminProtected := admin.Group("", middleware.AuthMiddleware(adminUsecase))
	adminProtected.Post("/logout", adminHandler.Logout)
	// Registration creates a full-privilege admin account, so it must only be
	// reachable by an already-authenticated admin — never from the public
	// internet. The first admin is provisioned at startup from environment
	// variables (see bootstrap.ensureInitialAdmin).
	adminProtected.Post("/register", adminHandler.Register)
	// adminProtected.Get("/profile", adminHandler.GetProfile)
	// adminProtected.Put("/profile", adminHandler.UpdateProfile)
	// adminProtected.Post("/change-password", adminHandler.ChangePassword)
}
