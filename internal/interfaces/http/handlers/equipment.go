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

// GetByID returns equipment detail by ID code
// GET /api/equipment/:id
func (h *EquipmentHandler) GetByID(c *fiber.Ctx) error {
	idCode := c.Params("id")
	if idCode == "" {
		return errors.BadRequest(c, "Equipment ID is required")
	}

	result, err := h.equipmentUsecase.GetByIDCode(c.Context(), idCode)
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, result, "Equipment retrieved successfully")
}

// Update updates equipment by ID code
// PUT /api/equipment/:id
func (h *EquipmentHandler) Update(c *fiber.Ctx) error {
	idCode := c.Params("id")
	if idCode == "" {
		return errors.BadRequest(c, "Equipment ID is required")
	}

	var req dto.EquipmentUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest(c, "Invalid request body")
	}

	if err := h.equipmentUsecase.UpdateEquipment(c.Context(), idCode, req); err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, nil, "Equipment updated successfully")
}

// Delete soft deletes equipment by ID code
// DELETE /api/equipment/:id
func (h *EquipmentHandler) Delete(c *fiber.Ctx) error {
	idCode := c.Params("id")
	if idCode == "" {
		return errors.BadRequest(c, "Equipment ID is required")
	}

	if err := h.equipmentUsecase.DeleteEquipment(c.Context(), idCode); err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, nil, "Equipment deleted successfully")
}
