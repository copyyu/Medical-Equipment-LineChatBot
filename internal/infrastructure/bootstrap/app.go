package bootstrap

import (
	"fmt"
	"log"
	"medical-webhook/internal/application/mapper"
	"medical-webhook/internal/application/service"
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/config"
	"medical-webhook/internal/domain/event"
	"medical-webhook/internal/infrastructure/client"
	"medical-webhook/internal/infrastructure/database"
	"medical-webhook/internal/infrastructure/logger"
	"medical-webhook/internal/infrastructure/persistence"
	redisinfra "medical-webhook/internal/infrastructure/redis"
	"medical-webhook/internal/infrastructure/session"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"
	"medical-webhook/internal/interfaces/http/routes"
	apperrors "medical-webhook/internal/utils/errors"
	"medical-webhook/internal/utils/scheduler"

	"github.com/gofiber/fiber/v2"
)

type Application struct {
	Server                 *fiber.App
	Config                 *config.Config
	WebhookHandler         *handlers.WebhookHandler
	NotificationHandler    *handlers.NotificationHandler
	NotificationScheduler  *scheduler.NotificationScheduler
	EquipmentImportHandler *handlers.EquipmentImportHandler
	AdminHandler           *handlers.AdminHandler
	EquipmentHandler       *handlers.EquipmentHandler
	TicketHandler          *handlers.TicketHandler
	ActivityLogHandler     *handlers.ActivityLogHandler
	SSEHandler             *handlers.SSEHandler
}

// InitializeApp - setup dependencies, routes, and return ready-to-run Application
func InitializeApp() (*Application, func(), error) {
	// Load and validate configuration (fail fast on missing required env vars,
	// e.g. an empty LINE_CHANNEL_SECRET that would make webhook signatures forgeable)
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize structured logging (bridges the standard log package too)
	logger.Init(cfg.AppEnv, cfg.LogLevel)
	log.Printf("Starting application (env=%s, log_level=%s)", cfg.AppEnv, cfg.LogLevel)

	// Connect Database
	if err := database.Connect(cfg); err != nil {
		return nil, nil, err
	}

	// Create the first super-admin if none exists (idempotent)
	database.SeedBootstrapAdmin(database.GetDB(), cfg.Admin.Username, cfg.Admin.Password, cfg.Admin.Email, cfg.Admin.FullName)

	// Initialize LINE client
	lineClient, err := client.NewClient(cfg.LineChannelToken, cfg.HTTP.LineAPITimeout)
	if err != nil {
		return nil, nil, err
	}

	// Initialize OCR client (optional - may be nil if not configured)
	var ocrClient *client.OCRClient
	if cfg.OCRURL != "" {
		ocrClient = client.NewOCRClient(cfg.OCRURL, cfg.HTTP.OCRAPITimeout)
		log.Printf("OCR client initialized: %s", cfg.OCRURL)
	} else {
		log.Println("OCR_API_URL not configured, OCR features disabled")
	}

	// Connect Redis (for Event Bus / Pub/Sub)
	if err := redisinfra.Connect(cfg.RedisURL); err != nil {
		log.Printf("⚠️ Redis connection failed: %v (real-time events disabled)", err)
	} else {
		log.Println("✅ Redis connected for Event Bus")
	}

	// Create EventBus (nil-safe — using interface type so nil stays nil)
	var eventBus event.EventBus
	if redisinfra.GetClient() != nil {
		eventBus = redisinfra.NewEventBus(redisinfra.GetClient())
	}

	// Idempotency store for de-duplicating LINE webhook events (nil if Redis down)
	var idempotencyStore *redisinfra.IdempotencyStore
	if redisinfra.GetClient() != nil {
		idempotencyStore = redisinfra.NewIdempotencyStore(redisinfra.GetClient(), 0)
	}

	// Initialize repositories (Infrastructure Layer)
	lineRepo := persistence.NewLineRepository(lineClient)
	notificationRepo := persistence.NewNotificationRepository(database.GetDB())
	equipmentRepo := persistence.NewEquipmentRepository()
	brandRepo := persistence.NewBrandRepository()
	equipmentCategoryRepo := persistence.NewEquipmentCategoryRepository()
	departmentRepo := persistence.NewDepartmentRepository()
	equipmentModelRepo := persistence.NewEquipmentModelRepository()
	adminRepo := persistence.NewAdminRepository()
	adminSessionRepo := persistence.NewAdminSessionRepository()

	// Initialize session store for OCR confirmations
	sessionStore := session.NewSessionStore()

	equipmentMapper := mapper.NewEquipmentMapper()

	// Initialize services (Domain Layer)
	messageService := service.NewMessageService(cfg.Contact)
	notificationService := service.NewNotificationService()
	excelParserService := service.NewExcelParserService()
	masterDataService := service.NewMasterDataService(
		brandRepo,
		equipmentCategoryRepo,
		departmentRepo,
		equipmentModelRepo,
		equipmentMapper,
	)

	adminService := service.NewAdminService(
		adminRepo,
		adminSessionRepo,
	)

	equipmentService := service.NewEquipmentService(
		equipmentRepo,
		brandRepo,
		equipmentCategoryRepo,
		departmentRepo,
		equipmentModelRepo,
	)

	// Initialize transaction manager
	txManager := database.NewTxManager(database.GetDB())

	// Initialize use cases (Application Layer)
	ticketRepo := persistence.NewTicketRepository(database.GetDB())
	ticketCategoryRepo := persistence.NewTicketCategoryRepository(database.GetDB())
	ticketHistoryRepo := persistence.NewTicketHistoryRepository(database.GetDB())
	ticketNotifyService := service.NewTicketNotificationService(lineRepo, ticketRepo)
	ticketUseCase := usecase.NewTicketUseCase(
		lineRepo,
		equipmentRepo,
		ticketRepo,
		ticketCategoryRepo,
		ticketHistoryRepo,
		ticketNotifyService,
		eventBus,
		txManager,
	)

	messageUseCase := usecase.NewMessageUseCase(
		lineRepo,
		equipmentRepo,
		departmentRepo,
		ocrClient,
		sessionStore,
		messageService,
		ticketUseCase,
		cfg.BaseURL,
	)
	notificationUseCase := usecase.NewNotificationUseCase(
		notificationRepo,
		notificationService,
		lineRepo,
		equipmentRepo,
	)

	equipmentImportUseCase := usecase.NewEquipmentImportUseCase(
		equipmentRepo,
		excelParserService,
		masterDataService,
		equipmentMapper,
	)

	adminUseCase := usecase.NewAdminUsecase(
		adminService,
	)

	// Initialize maintenance repository for dashboard
	maintenanceRepo := persistence.NewMaintenanceRecordRepository()

	dashboardUseCase := usecase.NewDashboardUsecase(
		equipmentRepo,
		maintenanceRepo,
		ticketRepo,
	)

	// Initialize equipment usecase for equipment list (using service layer)
	equipmentUseCase := usecase.NewEquipmentUsecase(equipmentService, eventBus)

	// Initialize activity log usecase (reuses ticketHistoryRepo)
	activityLogUseCase := usecase.NewActivityLogUseCase(ticketHistoryRepo)

	// Initialize handlers (Interface Layer)
	webhookHandler := handlers.NewWebhookHandler(cfg.LineChannelSecret, messageUseCase, idempotencyStore)
	notificationHandler := handlers.NewNotificationHandler(notificationUseCase)
	equipmentImportHandler := handlers.NewEquipmentImportHandler(equipmentImportUseCase)
	adminHandler := handlers.NewAdminHandler(adminUseCase)
	dashboardHandler := handlers.NewDashboardHandler(dashboardUseCase)
	equipmentHandler := handlers.NewEquipmentHandler(equipmentUseCase)
	ticketHandler := handlers.NewTicketHandler(ticketUseCase)
	activityLogHandler := handlers.NewActivityLogHandler(activityLogUseCase)

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		AppName:      "Medical Equipment Webhook",
		ErrorHandler: apperrors.FiberErrorHandler,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	})

	// Initialize SSE handler for real-time event streaming
	// Always create handler — it returns 503 gracefully if Redis is not connected
	sseHandler := handlers.NewSSEHandler(eventBus)

	// Register Middlewares
	middleware.FiberMiddleware(app, cfg.HTTP.AllowedOrigins)

	// Register Routes (SSE handler passed for public registration before 404 catch-all)
	routes.Setup(app, webhookHandler, notificationHandler, equipmentImportHandler, adminHandler, dashboardHandler, equipmentHandler, ticketHandler, activityLogHandler, sseHandler, adminUseCase)

	log.Println("📡 SSE endpoint registered: /api/events/stream")

	// Initialize และ Start Notification Scheduler
	notificationScheduler := scheduler.NewNotificationScheduler(notificationUseCase)
	notificationScheduler.Start()
	log.Println("Notification scheduler started")

	// Cleanup function
	cleanup := func() {
		log.Println("Shutting down gracefully...")
		// Stop scheduler
		if notificationScheduler != nil {
			notificationScheduler.Stop()
		}
		// Close event bus
		if eventBus != nil {
			eventBus.Close()
		}
		// Close Redis
		if err := redisinfra.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
		// Close session store (stops cleanup goroutine)
		if sessionStore != nil {
			sessionStore.Close()
		}
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		log.Println("Cleanup complete")
	}

	return &Application{
		Server:                 app,
		Config:                 cfg,
		WebhookHandler:         webhookHandler,
		NotificationHandler:    notificationHandler,
		EquipmentImportHandler: equipmentImportHandler,
		AdminHandler:           adminHandler,
		EquipmentHandler:       equipmentHandler,
		TicketHandler:          ticketHandler,
		ActivityLogHandler:     activityLogHandler,
		SSEHandler:             sseHandler,
	}, cleanup, nil
}

// Start - start the server
func (a *Application) Start() error {
	return a.Server.Listen(":" + a.Config.Port)
}

// Shutdown - graceful shutdown, bounded by the configured timeout so a stuck
// in-flight request cannot block the process from exiting indefinitely.
func (a *Application) Shutdown() error {
	return a.Server.ShutdownWithTimeout(a.Config.HTTP.ShutdownTimeout)
}
