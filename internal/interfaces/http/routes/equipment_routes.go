package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupEquipmentRoutes(app *fiber.App, equipmentHandler *handlers.EquipmentHandler, adminUsecase usecase.AdminUsecase) {
	api := app.Group("/api")

	// Equipment routes - protected
	equipment := api.Group("/equipment", middleware.AuthMiddleware(adminUsecase))

	// GET /api/equipment - Get paginated equipment list
	equipment.Get("/", equipmentHandler.GetList)

	// GET /api/equipment/categories - Get all equipment categories
	equipment.Get("/categories", equipmentHandler.GetCategories)

	// POST /api/equipment - Create new equipment
	equipment.Post("/", equipmentHandler.CreateEquipment)

	// GET /api/equipment/:id - Get equipment by ID code
	equipment.Get("/:id", equipmentHandler.GetByID)

	// PUT /api/equipment/:id - Update equipment
	equipment.Put("/:id", equipmentHandler.Update)

	// DELETE /api/equipment/:id - Delete equipment
	equipment.Delete("/:id", equipmentHandler.Delete)
}
