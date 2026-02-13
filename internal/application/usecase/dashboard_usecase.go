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
	ticketRepo      repository.TicketRepository
}

func NewDashboardUsecase(
	equipmentRepo repository.EquipmentRepository,
	maintenanceRepo repository.MaintenanceRecordRepository,
	ticketRepo repository.TicketRepository,
) DashboardUsecase {
	return &dashboardUsecase{
		equipmentRepo:   equipmentRepo,
		maintenanceRepo: maintenanceRepo,
		ticketRepo:      ticketRepo,
	}
}

func (u *dashboardUsecase) GetDashboardSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error) {
	// 1. Get total equipment count
	totalEquipment, err := u.equipmentRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// ✅ 2. Get expired equipment count (remain_life <= 0)
	// เพิ่มส่วนนี้เพื่อดึงข้อมูลอุปกรณ์ที่หมดอายุจาก Database
	expiredEquipment, err := u.equipmentRepo.CountExpired(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Get near expiry count (equipment with remain_life <= 1 year)
	nearExpiry, err := u.equipmentRepo.CountNearExpiry(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Get total maintenance count
	totalMaintenance, err := u.maintenanceRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 5. Get equipment count by Asset Status
	assetStatusMap, err := u.equipmentRepo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	// 6. Get ticket stats (using TicketStatus)
	_, inProgress, completed, sendToOutsource, err := u.ticketRepo.GetTicketStats()
	if err != nil {
		return nil, err
	}

	// 7. Get recent maintenance jobs (from Tickets)
	recentTickets, err := u.ticketRepo.GetRecentTickets(5)
	if err != nil {
		return nil, err
	}

	// Build Asset Status Counts
	assetStatusCounts := []dto.AssetStatusCount{
		{Status: "active", Count: assetStatusMap[entity.AssetStatusActive]},
		{Status: "defective", Count: assetStatusMap[entity.AssetStatusDefective]},
		{Status: "wait_decom", Count: assetStatusMap[entity.AssetStatusWaitDecom]},
		{Status: "decommission", Count: assetStatusMap[entity.AssetStatusDecommission]},
		{Status: "active_ready_to_sell", Count: assetStatusMap[entity.AssetStatusActiveReadyToSell]},
		{Status: "missing", Count: assetStatusMap[entity.AssetStatusMissing]},
		{Status: "plan_to_replace", Count: assetStatusMap[entity.AssetStatusPlanToReplace]},
	}

	// Build Ticket Status Counts
	ticketStatusCounts := []dto.TicketStatusCount{
		{Status: "in_progress", Count: inProgress},
		{Status: "return_equipment_back", Count: completed},
		{Status: "send_to_outsource", Count: sendToOutsource},
	}

	// Build recent jobs list
	recentJobs := make([]dto.RecentJobResponse, 0)
	for _, t := range recentTickets {
		equipmentName := ""
		if t.Equipment != nil {
			if t.Equipment.Model.ModelName != "" {
				equipmentName = t.Equipment.Model.ModelName + " #" + t.Equipment.IDCode
			} else {
				equipmentName = t.Equipment.IDCode
			}
		} else {
			equipmentName = "Unknown Equipment"
		}

		assignee := "ยังไม่ได้มอบหมาย"

		recentJobs = append(recentJobs, dto.RecentJobResponse{
			ID:            t.TicketNo, // Use TicketNo instead of generated JOB-ID
			EquipmentName: equipmentName,
			Status:        string(t.Status), // Use TicketStatus
			Assignee:      assignee,
			UpdatedAt:     formatTimeAgo(t.UpdatedAt),
		})
	}

	return &dto.DashboardSummaryResponse{
		TotalEquipment:     totalEquipment,
		RentalEquipment:    expiredEquipment, // ✅ ใส่ค่า expiredEquipment ลงไปใน field นี้เพื่อให้ Frontend แสดงผล
		NearExpiry:         nearExpiry,
		TotalMaintenance:   totalMaintenance,
		AssetStatusCounts:  assetStatusCounts,
		TicketStatusCounts: ticketStatusCounts,
		RecentJobs:         recentJobs,
	}, nil
}

// Helper functions
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
