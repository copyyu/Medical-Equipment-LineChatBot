package usecase

import (
	"context"
	"fmt"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	"time"
)

type DashboardUsecase interface {
	GetDashboardSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error)
}

type dashboardUsecase struct {
	equipmentRepo   repository.EquipmentRepository
	maintenanceRepo repository.MaintenanceRecordRepository
}

func NewDashboardUsecase(
	equipmentRepo repository.EquipmentRepository,
	maintenanceRepo repository.MaintenanceRecordRepository,
) DashboardUsecase {
	return &dashboardUsecase{
		equipmentRepo:   equipmentRepo,
		maintenanceRepo: maintenanceRepo,
	}
}

func (u *dashboardUsecase) GetDashboardSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error) {
	// Get total equipment count
	totalEquipment, err := u.equipmentRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Get near expiry count (equipment with remain_life <= 1 year)
	nearExpiry, err := u.equipmentRepo.CountNearExpiry(ctx)
	if err != nil {
		return nil, err
	}

	// Get total maintenance count
	totalMaintenance, err := u.maintenanceRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Get equipment count by Asset Status (จาก field status ใน equipments table)
	assetStatusMap, err := u.equipmentRepo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	// Get maintenance count by Job Status (จาก field status ใน maintenance_records table)
	jobStatusMap, err := u.maintenanceRepo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	// Get recent maintenance jobs
	recentMaintenance, err := u.maintenanceRepo.GetRecent(ctx, 5)
	if err != nil {
		return nil, err
	}

	// Build Asset Status Counts from real data
	assetStatusCounts := []dto.AssetStatusCount{
		{Status: "active", Count: assetStatusMap[entity.AssetStatusActive]},
		{Status: "defective", Count: assetStatusMap[entity.AssetStatusDefective]},
		{Status: "wait_decom", Count: assetStatusMap[entity.AssetStatusWaitDecom]},
		{Status: "decommission", Count: assetStatusMap[entity.AssetStatusDecommission]},
		{Status: "active_ready_to_sell", Count: assetStatusMap[entity.AssetStatusActiveReadyToSell]},
		{Status: "missing", Count: assetStatusMap[entity.AssetStatusMissing]},
		{Status: "plan_to_replace", Count: assetStatusMap[entity.AssetStatusPlanToReplace]},
	}

	// Build Job Status Counts from real data
	jobCounts := []dto.JobStatusCount{
		{Status: "in_process", Count: jobStatusMap[entity.JobStatusInProcess]},
		{Status: "return_equipment_back", Count: jobStatusMap[entity.JobStatusReturnEquipmentBack]},
		{Status: "send_to_outsource", Count: jobStatusMap[entity.JobStatusSendToOutsource]},
	}

	// Build recent jobs list
	recentJobs := make([]dto.RecentJobResponse, 0)
	for _, m := range recentMaintenance {
		equipmentName := ""
		if m.Equipment.Model.ModelName != "" {
			equipmentName = m.Equipment.Model.ModelName + " #" + m.Equipment.IDCode
		} else {
			equipmentName = m.Equipment.IDCode
		}

		recentJobs = append(recentJobs, dto.RecentJobResponse{
			ID:            formatJobID(m.ID),
			EquipmentName: equipmentName,
			Status:        string(m.Status), // ใช้ status จริงจาก DB
			Assignee:      getAssignee(m.Technician),
			UpdatedAt:     formatTimeAgo(m.UpdatedAt),
		})
	}

	return &dto.DashboardSummaryResponse{
		TotalEquipment:    totalEquipment,
		RentalEquipment:   0, // TODO: Add rental tracking if needed
		NearExpiry:        nearExpiry,
		TotalMaintenance:  totalMaintenance,
		AssetStatusCounts: assetStatusCounts,
		JobStatusCounts:   jobCounts,
		RecentJobs:        recentJobs,
	}, nil
}

func formatJobID(id uint) string {
	return fmt.Sprintf("JOB-2026-%04d", id)
}

func getAssignee(name string) string {
	if name == "" {
		return "ยังไม่ได้มอบหมาย"
	}
	return name
}

func formatTimeAgo(t time.Time) string {
	diff := time.Since(t)

	if diff < time.Minute {
		return "เมื่อสักครู่"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())
		return fmt.Sprintf("%d นาทีที่แล้ว", mins)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d ชั่วโมงที่แล้ว", hours)
	} else {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d วันที่แล้ว", days)
	}
}
