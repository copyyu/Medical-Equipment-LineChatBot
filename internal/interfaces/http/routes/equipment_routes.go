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
	equipment.Post("/", equipmentHandler.CreateEquipment)
}
