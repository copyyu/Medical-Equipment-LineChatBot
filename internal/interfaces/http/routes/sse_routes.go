package routes

import (
	"medical-webhook/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupSSERoutes registers the Server-Sent Events routes
func SetupSSERoutes(app *fiber.App, sseHandler *handlers.SSEHandler) {
	// SSE endpoint is public — no auth required
	// Frontend opens EventSource connection to receive real-time updates
	app.Get("/api/events/stream", sseHandler.Stream)
}
