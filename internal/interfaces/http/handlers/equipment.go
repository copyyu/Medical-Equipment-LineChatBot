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

// GetByID returns equipment by ID
// GET /api/equipment/:id
func (h *EquipmentHandler) GetByID(c *fiber.Ctx) error {
	log.Printf("Handler: GetByID - Received request")

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Printf("Handler: GetByID - Invalid ID parameter: %v", err)
		return errors.Error(c, fiber.NewError(fiber.StatusBadRequest, "Invalid equipment ID"))
	}

	log.Printf("Handler: GetByID - ID: %d", id)

	return errors.Error(c, fiber.NewError(fiber.StatusNotImplemented, "GetByID not implemented yet"))
}

// UpdateEquipment updates an existing equipment
// PUT /api/equipment/:id
func (h *EquipmentHandler) UpdateEquipment(c *fiber.Ctx) error {
	log.Printf("Handler: UpdateEquipment - Received request")

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Printf("Handler: UpdateEquipment - Invalid ID parameter: %v", err)
		return errors.Error(c, fiber.NewError(fiber.StatusBadRequest, "Invalid equipment ID"))
	}

	var req dto.CreateEquipmentRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Handler: UpdateEquipment - Body parse error: %v", err)
		return errors.Error(c, fiber.NewError(fiber.StatusBadRequest, "Invalid request body"))
	}

	log.Printf("Handler: UpdateEquipment - ID: %d, IDCode: %s", id, req.IDCode)

	return errors.Error(c, fiber.NewError(fiber.StatusNotImplemented, "UpdateEquipment not implemented yet"))
}

// DeleteEquipment deletes an equipment
// DELETE /api/equipment/:id
func (h *EquipmentHandler) DeleteEquipment(c *fiber.Ctx) error {
	log.Printf("Handler: DeleteEquipment - Received request")

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Printf("Handler: DeleteEquipment - Invalid ID parameter: %v", err)
		return errors.Error(c, fiber.NewError(fiber.StatusBadRequest, "Invalid equipment ID"))
	}

	log.Printf("Handler: DeleteEquipment - ID: %d", id)

	return errors.Error(c, fiber.NewError(fiber.StatusNotImplemented, "DeleteEquipment not implemented yet"))
}
