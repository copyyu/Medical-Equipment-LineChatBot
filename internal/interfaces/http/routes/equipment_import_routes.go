package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupEquipmentImportRoutes
func SetupEquipmentImportRoutes(app *fiber.App, equipmentImportHandler *handlers.EquipmentImportHandler) {
	app.Post("/import", equipmentImportHandler.ImportExcel)
	app.Post("/batch", equipmentImportHandler.ImportExcelBatch)
	app.Get("/history", equipmentImportHandler.GetImportHistory)

}
