package usecase

import (
	"context"
	"errors"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	"time"
)

type EquipmentUsecase interface {
	GetEquipmentList(ctx context.Context, req dto.EquipmentListRequest) (*dto.EquipmentListResponse, error)
	GetByIDCode(ctx context.Context, idCode string) (*dto.EquipmentDetailResponse, error)
	UpdateEquipment(ctx context.Context, idCode string, req dto.EquipmentUpdateRequest) error
	DeleteEquipment(ctx context.Context, idCode string) error
}

type equipmentUsecase struct {
	equipmentRepo  repository.EquipmentRepository
	departmentRepo repository.DepartmentRepository
}

func NewEquipmentUsecase(equipmentRepo repository.EquipmentRepository, departmentRepo repository.DepartmentRepository) EquipmentUsecase {
	return &equipmentUsecase{
		equipmentRepo:  equipmentRepo,
		departmentRepo: departmentRepo,
	}
}

func (u *equipmentUsecase) GetEquipmentList(ctx context.Context, req dto.EquipmentListRequest) (*dto.EquipmentListResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit
	}

	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get total count with filters
	total, err := u.equipmentRepo.CountWithFilter(ctx, req.Status, req.Search)
	if err != nil {
		return nil, err
	}

	// Get equipment list with pagination and filters
	equipments, err := u.equipmentRepo.FindAllWithFilter(ctx, req.Limit, offset, req.Status, req.Search)
	if err != nil {
		return nil, err
	}

	// Map to DTO
	items := make([]dto.EquipmentListItem, 0, len(equipments))
	for _, e := range equipments {
		items = append(items, dto.MapEquipmentToListItem(e))
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.EquipmentListResponse{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetByIDCode returns equipment detail by ID code
func (u *equipmentUsecase) GetByIDCode(ctx context.Context, idCode string) (*dto.EquipmentDetailResponse, error) {
	equipment, err := u.equipmentRepo.FindByIDCode(idCode)
	if err != nil {
		return nil, err
	}
	if equipment == nil {
		return nil, errors.New("equipment not found")
	}

	result := dto.MapEquipmentToDetailResponse(*equipment)
	return &result, nil
}

// UpdateEquipment updates equipment by ID code
func (u *equipmentUsecase) UpdateEquipment(ctx context.Context, idCode string, req dto.EquipmentUpdateRequest) error {
	// Find existing equipment
	equipment, err := u.equipmentRepo.FindByIDCode(idCode)
	if err != nil {
		return err
	}
	if equipment == nil {
		return errors.New("equipment not found")
	}

	// Update status if provided
	if req.Status != "" {
		equipment.Status = entity.AssetStatus(req.Status)
	}

	// Update department if location provided
	if req.Location != "" {
		dept, err := u.departmentRepo.FindOrCreate(ctx, req.Location)
		if err != nil {
			return err
		}
		equipment.DepartmentID = dept.ID
	}

	// Update compute date if provided
	if req.ComputeDate != "" {
		computeDate, err := time.Parse("2006-01-02", req.ComputeDate)
		if err == nil {
			equipment.ComputeDate = &computeDate
		}
	}

	return u.equipmentRepo.Update(ctx, equipment)
}

// DeleteEquipment soft deletes equipment by ID code
func (u *equipmentUsecase) DeleteEquipment(ctx context.Context, idCode string) error {
	// Find existing equipment
	equipment, err := u.equipmentRepo.FindByIDCode(idCode)
	if err != nil {
		return err
	}
	if equipment == nil {
		return errors.New("equipment not found")
	}

	return u.equipmentRepo.Delete(ctx, equipment.ID)
}
