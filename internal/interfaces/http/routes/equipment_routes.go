package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupEquipmentRoutes(app *fiber.App, equipmentHandler *handlers.EquipmentHandler) {
	// API routes
	api := app.Group("/api")

	// Equipment routes
	equipment := api.Group("/equipment")

	// GET /api/equipment - Get paginated equipment list
	equipment.Get("/", equipmentHandler.GetList)

	// POST /api/equipment - Create new equipment
	equipment.Post("/", equipmentHandler.CreateEquipment)

	// GET /api/equipment/:id - Get equipment by ID code
	equipment.Get("/:id", equipmentHandler.GetByID)

	// PUT /api/equipment/:id - Update equipment
	equipment.Put("/:id", equipmentHandler.Update)

	// DELETE /api/equipment/:id - Delete equipment
	equipment.Delete("/:id", equipmentHandler.Delete)
}
