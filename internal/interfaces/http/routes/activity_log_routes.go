package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupActivityLogRoutes(app *fiber.App, activityLogHandler *handlers.ActivityLogHandler, adminUsecase usecase.AdminUsecase) {
	api := app.Group("/api")

	// Activity log routes - protected
	activityLog := api.Group("/activity-logs", middleware.AuthMiddleware(adminUsecase))

	// GET /api/activity-logs - Get paginated activity log list
	activityLog.Get("/", activityLogHandler.GetList)

	// GET /api/activity-logs/stats - Get activity log stats
	activityLog.Get("/stats", activityLogHandler.GetStats)
}
