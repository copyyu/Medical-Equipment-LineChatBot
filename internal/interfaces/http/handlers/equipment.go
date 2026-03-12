package handlers

import (
	"log"
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
		Page:         c.QueryInt("page", 1),
		Limit:        c.QueryInt("limit", 10),
		Status:       c.Query("status"),
		Search:       c.Query("search"),
		SortBy:       c.Query("sort_by", "id"),
		SortDir:      c.Query("sort_dir", "desc"),
		ExpiryFilter: c.Query("expiry_filter"),
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

// CreateEquipment creates a new equipment
// POST /api/equipment
// Body: CreateEquipmentRequest JSON
func (h *EquipmentHandler) CreateEquipment(c *fiber.Ctx) error {
	log.Printf("Handler: CreateEquipment - Received request")

	var req dto.CreateEquipmentRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Handler: CreateEquipment - Body parse error: %v", err)
		return errors.Error(c, fiber.NewError(fiber.StatusBadRequest, "Invalid request body"))
	}

	log.Printf("Handler: CreateEquipment - IDCode: %s, SerialNo: %s", req.IDCode, req.SerialNo)

	result, err := h.equipmentUsecase.CreateEquipment(c.Context(), req)
	if err != nil {
		log.Printf("Handler: CreateEquipment - Error: %v", err)
		return errors.Error(c, err)
	}

	log.Printf("Handler: CreateEquipment - Success, created equipment ID: %s", req.IDCode)
	return errors.Success(c, result, "Equipment created successfully")
}

// GetCategories returns all equipment categories
// GET /api/equipment/categories
func (h *EquipmentHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.equipmentUsecase.GetAllCategories(c.Context())
	if err != nil {
		log.Printf("Handler: GetCategories - Error: %v", err)
		return errors.Error(c, err)
	}

	return errors.Success(c, categories, "Categories retrieved successfully")
}
