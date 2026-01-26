package scheduler

import (
	"context"
	"log"

	"medical-webhook/internal/application/usecase"

	"github.com/robfig/cron/v3"
)

type NotificationScheduler struct {
	cron                *cron.Cron
	notificationUseCase *usecase.NotificationUseCase
}

func NewNotificationScheduler(notificationUseCase *usecase.NotificationUseCase) *NotificationScheduler {
	return &NotificationScheduler{
		cron:                cron.New(),
		notificationUseCase: notificationUseCase,
	}
}

// Start - เริ่ม scheduler
func (s *NotificationScheduler) Start() {
	// ส่งแจ้งเตือนรอบมิถุนายน: วันที่ 1 มิ.ย. เวลา 09:00
	s.cron.AddFunc("0 9 1 6 *", func() {
		log.Println("Running June notification...")
		ctx := context.Background()
		if err := s.notificationUseCase.SendJuneAlerts(ctx); err != nil {
			log.Printf("Error sending June alerts: %v", err)
		}
	})

	// ส่งแจ้งเตือนรอบสิงหาคม: วันที่ 1 ส.ค. เวลา 09:00
	s.cron.AddFunc("0 9 1 8 *", func() {
		log.Println("Running August notification...")
		ctx := context.Background()
		if err := s.notificationUseCase.SendAugustAlerts(ctx); err != nil {
			log.Printf("Error sending August alerts: %v", err)
		}
	})

	s.cron.Start()
	log.Println("Notification scheduler started")
	log.Println("Schedule: June 1st at 09:00 & August 1st at 09:00")
}

// Test Cronjob
// func (s *NotificationScheduler) Start() {
// 	// ทดสอบ - ส่งทุก 1 นาที
// 	s.cron.AddFunc("*/1 * * * *", func() {
// 		log.Println("[TEST] Running June notification...")
// 		ctx := context.Background()
// 		if err := s.notificationUseCase.SendJuneAlerts(ctx); err != nil {
// 			log.Printf("Error sending June alerts: %v", err)
// 		}
// 	})

// 	s.cron.Start()
// 	log.Println("TEST MODE: Notification scheduler running every 1 minute")
// }

// Stop - หยุด scheduler
func (s *NotificationScheduler) Stop() {
	s.cron.Stop()
	log.Println("Notification scheduler stopped")
}
