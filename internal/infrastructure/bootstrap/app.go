package bootstrap

import (
	"log"
	"medical-webhook/internal/application/mapper"
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/config"
	"medical-webhook/internal/domain/line/service"
	"medical-webhook/internal/infrastructure/client"
	"medical-webhook/internal/infrastructure/database"
	"medical-webhook/internal/infrastructure/persistence"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"
	"medical-webhook/internal/interfaces/http/routes"
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
}

// InitializeApp - setup dependencies, routes, and return ready-to-run Application
func InitializeApp() (*Application, func(), error) {
	// Load configuration
	cfg := config.Load()

	// Connect Database
	if err := database.Connect(cfg); err != nil {
		return nil, nil, err
	}

	// Initialize LINE client
	lineClient, err := client.NewClient(cfg.LineChannelToken)
	if err != nil {
		return nil, nil, err
	}

	// Initialize OCR client (optional - may be nil if not configured)
	var ocrClient *client.OCRClient
	if cfg.OCRURL != "" {
		ocrClient = client.NewOCRClient(cfg.OCRURL)
		log.Printf("OCR client initialized: %s", cfg.OCRURL)
	} else {
		log.Println("OCR_API_URL not configured, OCR features disabled")
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
	sessionStore := usecase.NewSessionStore()

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

	// Initialize use cases (Application Layer)
	messageUseCase := usecase.NewMessageUseCase(
		lineRepo,
		equipmentRepo,
		ocrClient,
		sessionStore,
		messageService,
	)
	notificationUseCase := usecase.NewNotificationUseCase(
		notificationRepo,
		notificationService,
		lineRepo,
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
	)

	// Initialize equipment usecase for equipment list
	equipmentUseCase := usecase.NewEquipmentUsecase(equipmentRepo, departmentRepo)

	// Initialize handlers (Interface Layer)
	webhookHandler := handlers.NewWebhookHandler(cfg.LineChannelSecret, messageUseCase)
	notificationHandler := handlers.NewNotificationHandler(notificationUseCase)
	equipmentImportHandler := handlers.NewEquipmentImportHandler(equipmentImportUseCase)
	adminHandler := handlers.NewAdminHandler(adminUseCase)
	dashboardHandler := handlers.NewDashboardHandler(dashboardUseCase)
	equipmentHandler := handlers.NewEquipmentHandler(equipmentUseCase)

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		AppName: "Medical Equipment Webhook",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		},
	})

	// Register Middlewares
	middleware.FiberMiddleware(app)

	// Register Routes
	routes.Setup(app, webhookHandler, notificationHandler, equipmentImportHandler, adminHandler, dashboardHandler, equipmentHandler)

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
	}, cleanup, nil
}

// Start - start the server
func (a *Application) Start() error {
	return a.Server.Listen(":" + a.Config.Port)
}

// Shutdown - graceful shutdown
func (a *Application) Shutdown() error {
	return a.Server.Shutdown()
}
