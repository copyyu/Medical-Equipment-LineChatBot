package routes

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTicketRoutes(app *fiber.App, ticketHandler *handlers.TicketHandler, adminUsecase usecase.AdminUsecase) {
	api := app.Group("/api")

	// Ticket routes - protected
	ticket := api.Group("/tickets", middleware.AuthMiddleware(adminUsecase))

	// GET /api/tickets - Get paginated ticket list
	ticket.Get("/", ticketHandler.GetList)

	// GET /api/tickets/stats - Get ticket stats
	ticket.Get("/stats", ticketHandler.GetStats)

	// GET /api/tickets/categories - Get ticket categories
	ticket.Get("/categories", ticketHandler.GetCategories)

	// GET /api/tickets/:id - Get ticket detail
	ticket.Get("/:id", ticketHandler.GetByID)

	// PUT /api/tickets/:id - Update ticket details
	ticket.Put("/:id", ticketHandler.UpdateTicket)
}
