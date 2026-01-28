package handlers

import (
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
)

type EquipmentHandler struct {
	equipmentUsecase usecase.EquipmentUsecase
}

func NewEquipmentHandler(equipmentUsecase usecase.EquipmentUsecase) *EquipmentHandler {
	return &EquipmentHandler{
		equipmentUsecase: equipmentUsecase,
	}
}

// GetList returns paginated equipment list
// GET /api/equipment
// Query params: page, limit, status, search, sort_by, sort_dir
func (h *EquipmentHandler) GetList(c *fiber.Ctx) error {
	req := dto.EquipmentListRequest{
		Page:    c.QueryInt("page", 1),
		Limit:   c.QueryInt("limit", 10),
		Status:  c.Query("status"),
		Search:  c.Query("search"),
		SortBy:  c.Query("sort_by", "id"),
		SortDir: c.Query("sort_dir", "desc"),
	}

	result, err := h.equipmentUsecase.GetEquipmentList(c.Context(), req)
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, result, "Equipment list retrieved successfully")
}
