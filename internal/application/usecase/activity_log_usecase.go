package usecase

import (
	"context"
	"math"

	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/repository"
)

// ActivityLogUseCase handles activity log business logic
type ActivityLogUseCase struct {
	historyRepo repository.TicketHistoryRepository
}

// NewActivityLogUseCase creates a new activity log use case
func NewActivityLogUseCase(historyRepo repository.TicketHistoryRepository) *ActivityLogUseCase {
	return &ActivityLogUseCase{
		historyRepo: historyRepo,
	}
}

// GetActivityLogs returns paginated & filtered activity logs
func (uc *ActivityLogUseCase) GetActivityLogs(ctx context.Context, req dto.ActivityLogListRequest) (*dto.ActivityLogListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	entries, total, err := uc.historyRepo.GetStatusChangeLogs(
		req.Page, req.Limit,
		req.Search, req.FromStatus, req.ToStatus,
		req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	var items []dto.ActivityLogItem
	for _, e := range entries {
		items = append(items, dto.ActivityLogItem{
			ID:            e.ID,
			TicketID:      e.TicketID,
			TicketNo:      e.TicketNo,
			EquipmentName: e.EquipmentName,
			UserName:      e.AdminName,
			FromStatus:    e.FromStatus,
			ToStatus:      e.ToStatus,
			Note:          e.Note,
			ChangedAt:     e.ChangedAt,
		})
	}

	return &dto.ActivityLogListResponse{
		Data: items,
		Pagination: dto.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetActivityLogStats returns aggregate stats
func (uc *ActivityLogUseCase) GetActivityLogStats(ctx context.Context) (*dto.ActivityLogStatsResponse, error) {
	stats, err := uc.historyRepo.GetStatusChangeLogStats()
	if err != nil {
		return nil, err
	}

	return &dto.ActivityLogStatsResponse{
		TotalChanges:      stats.TotalChanges,
		TodayChanges:      stats.TodayChanges,
		WeekChanges:       stats.WeekChanges,
		MostChangedStatus: stats.MostChangedStatus,
	}, nil
}
