package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupEquipmentImportRoutes - protected routes for import
func SetupEquipmentImportRoutes(app *fiber.App, equipmentImportHandler *handlers.EquipmentImportHandler, adminUsecase usecase.AdminUsecase) {
	// Import routes - protected
	importGroup := app.Group("", middleware.AuthMiddleware(adminUsecase))
	importGroup.Post("/import", equipmentImportHandler.ImportExcel)
	importGroup.Post("/batch", equipmentImportHandler.ImportExcelBatch)
	importGroup.Get("/history", equipmentImportHandler.GetImportHistory)
}
