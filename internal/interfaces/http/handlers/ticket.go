package handlers

import (
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/utils/errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	ticketUseCase *usecase.TicketUseCase
}

func NewTicketHandler(ticketUseCase *usecase.TicketUseCase) *TicketHandler {
	return &TicketHandler{
		ticketUseCase: ticketUseCase,
	}
}

// GET /api/tickets
func (h *TicketHandler) GetList(c *fiber.Ctx) error {
	req := dto.TicketListRequest{
		Page:     c.QueryInt("page", 1),
		Limit:    c.QueryInt("limit", 10),
		Status:   c.Query("status"),
		Priority: c.Query("priority"),
		Search:   c.Query("search"),
		SortBy:   c.Query("sort_by", "created_at"),
		SortDir:  c.Query("sort_dir", "desc"),
	}

	result, err := h.ticketUseCase.GetTicketList(c.Context(), req)
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, result, "Ticket list retrieved successfully")
}

// GET /api/tickets/:id
func (h *TicketHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest(c, "Invalid ticket ID")
	}

	result, err := h.ticketUseCase.GetTicketByID(c.Context(), uint(id))
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, result, "Ticket retrieved successfully")
}

// PUT /api/tickets/:id
func (h *TicketHandler) UpdateTicket(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest(c, "Invalid ticket ID")
	}

	var req dto.UpdateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest(c, "Invalid request body")
	}

	if err := h.ticketUseCase.UpdateTicket(c.Context(), uint(id), req); err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, nil, "Ticket updated successfully")
}

// GET /api/tickets/stats
func (h *TicketHandler) GetStats(c *fiber.Ctx) error {
	result, err := h.ticketUseCase.GetTicketStats(c.Context())
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, result, "Ticket stats retrieved successfully")
}

// GET /api/tickets/categories
func (h *TicketHandler) GetCategories(c *fiber.Ctx) error {
	result, err := h.ticketUseCase.GetTicketCategories(c.Context())
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, result, "Ticket categories retrieved successfully")
}
