package handlers

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	dashboardUsecase usecase.DashboardUsecase
}

func NewDashboardHandler(dashboardUsecase usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{
		dashboardUsecase: dashboardUsecase,
	}
}

// GetSummary returns dashboard summary data
// GET /api/dashboard/summary
func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	summary, err := h.dashboardUsecase.GetDashboardSummary(c.Context())
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, summary, "Dashboard summary retrieved successfully")
}
