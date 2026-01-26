package usecase

import (
	"context"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/repository"
)

type EquipmentUsecase interface {
	GetEquipmentList(ctx context.Context, req dto.EquipmentListRequest) (*dto.EquipmentListResponse, error)
}

type equipmentUsecase struct {
	equipmentRepo repository.EquipmentRepository
}

func NewEquipmentUsecase(equipmentRepo repository.EquipmentRepository) EquipmentUsecase {
	return &equipmentUsecase{
		equipmentRepo: equipmentRepo,
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

	// Get total count
	total, err := u.equipmentRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Get equipment list with pagination
	equipments, err := u.equipmentRepo.FindAll(ctx, req.Limit, offset)
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
