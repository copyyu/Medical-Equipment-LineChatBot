// interfaces/http/routes/routes.go
package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App, adminHandler *handlers.AdminHandler) {
	// API routes
	api := app.Group("/api")

	// Admin routes
	admin := api.Group("/admin")

	// Public routes - ไม่ต้อง auth
	admin.Post("/register", adminHandler.Register)
	admin.Post("/login", adminHandler.Login)

	// Protected routes - ต้อง auth
	// ใช้ middleware.AuthMiddleware() สำหรับ routes group
	// adminProtected := admin.Group("", middleware.AuthMiddleware(adminUsecase))
	// adminProtected.Post("/logout", adminHandler.Logout)
	// adminProtected.Get("/profile", adminHandler.GetProfile)
	// adminProtected.Put("/profile", adminHandler.UpdateProfile)
	// adminProtected.Post("/change-password", adminHandler.ChangePassword)

	// ตัวอย่างการใช้ auth กับ route เดียว
	// admin.Get("/some-route", middleware.AuthMiddleware(adminUsecase), someHandler)
}
