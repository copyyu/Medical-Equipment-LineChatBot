package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupSSERoutes registers the Server-Sent Events routes. The stream carries
// real-time ticket/equipment data (reporter names, serials, ticket numbers), so
// it is authenticated. Because EventSource cannot set headers, the frontend must
// pass the token as a query param:
//
//	new EventSource(`/api/events/stream?token=${bearerToken}`)
func SetupSSERoutes(app *fiber.App, sseHandler *handlers.SSEHandler, adminUsecase usecase.AdminUsecase) {
	app.Get("/api/events/stream", middleware.SSEAuthMiddleware(adminUsecase), sseHandler.Stream)
}
