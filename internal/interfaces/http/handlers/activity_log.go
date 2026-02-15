package handlers

import (
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
)

type ActivityLogHandler struct {
	activityLogUseCase *usecase.ActivityLogUseCase
}

func NewActivityLogHandler(activityLogUseCase *usecase.ActivityLogUseCase) *ActivityLogHandler {
	return &ActivityLogHandler{
		activityLogUseCase: activityLogUseCase,
	}
}

// GET /api/activity-logs
func (h *ActivityLogHandler) GetList(c *fiber.Ctx) error {
	req := dto.ActivityLogListRequest{
		Page:       c.QueryInt("page", 1),
		Limit:      c.QueryInt("limit", 20),
		Search:     c.Query("search"),
		FromStatus: c.Query("from_status"),
		ToStatus:   c.Query("to_status"),
		StartDate:  c.Query("start_date"),
		EndDate:    c.Query("end_date"),
	}

	result, err := h.activityLogUseCase.GetActivityLogs(c.Context(), req)
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, result, "Activity logs retrieved successfully")
}

// GET /api/activity-logs/stats
func (h *ActivityLogHandler) GetStats(c *fiber.Ctx) error {
	result, err := h.activityLogUseCase.GetActivityLogStats(c.Context())
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, result, "Activity log stats retrieved successfully")
}
