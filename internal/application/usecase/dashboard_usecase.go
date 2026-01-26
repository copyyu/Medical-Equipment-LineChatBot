package usecase

import (
	"context"
	"fmt"
	"medical-webhook/internal/application/dto"
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

	// Get maintenance by type (CM/PM)
	maintenanceTypeCounts, err := u.maintenanceRepo.CountByType(ctx)
	if err != nil {
		return nil, err
	}

	// Get recent maintenance jobs
	recentMaintenance, err := u.maintenanceRepo.GetRecent(ctx, 5)
	if err != nil {
		return nil, err
	}

	// Map maintenance types to job status format for frontend
	// CM (Corrective) -> in_process, PM (Preventive) -> return_equipment_back
	jobCounts := []dto.JobStatusCount{
		{Status: "in_process", Count: maintenanceTypeCounts["CM"]},
		{Status: "return_equipment_back", Count: maintenanceTypeCounts["PM"]},
		{Status: "send_to_outsource", Count: 0},
	}

	recentJobs := make([]dto.RecentJobResponse, 0)
	for _, m := range recentMaintenance {
		equipmentName := ""
		if m.Equipment.Model.ModelName != "" {
			equipmentName = m.Equipment.Model.ModelName + " #" + m.Equipment.IDCode
		} else {
			equipmentName = m.Equipment.IDCode
		}

		// Map maintenance type to status
		status := "in_process"
		if m.MaintenanceType == "PM" {
			status = "return_equipment_back"
		}

		recentJobs = append(recentJobs, dto.RecentJobResponse{
			ID:            formatJobID(m.ID),
			EquipmentName: equipmentName,
			Status:        status,
			Assignee:      getAssignee(m.Technician),
			UpdatedAt:     formatTimeAgo(m.UpdatedAt),
		})
	}

	// Asset status counts based on equipment data
	assetStatusCounts := []dto.AssetStatusCount{
		{Status: "active", Count: totalEquipment - nearExpiry},
		{Status: "defective", Count: 0},
		{Status: "wait_decom", Count: 0},
		{Status: "decommission", Count: 0},
		{Status: "active_ready_to_sell", Count: 0},
		{Status: "missing", Count: 0},
		{Status: "plan_to_replace", Count: nearExpiry},
	}

	return &dto.DashboardSummaryResponse{
		TotalEquipment:    totalEquipment,
		RentalEquipment:   0, // No rental data in current schema
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
