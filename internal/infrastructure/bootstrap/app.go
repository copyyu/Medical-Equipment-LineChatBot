package bootstrap

import (
	"log"
	"medical-webhook/internal/application/mapper"
	"medical-webhook/internal/application/service"
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/config"
	"medical-webhook/internal/infrastructure/client"
	"medical-webhook/internal/infrastructure/database"
	"medical-webhook/internal/infrastructure/persistence"
	"medical-webhook/internal/infrastructure/session"
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
	EquipmentHandler       *handlers.EquipmentHandler
	TicketHandler          *handlers.TicketHandler
}

// repositories holds all infrastructure repository instances
type repositories struct {
	line              *persistence.LineRepository
	notification      *persistence.NotificationRepository
	equipment         *persistence.EquipmentRepository
	brand             *persistence.BrandRepository
	equipmentCategory *persistence.EquipmentCategoryRepository
	department        *persistence.DepartmentRepository
	equipmentModel    *persistence.EquipmentModelRepository
	admin             *persistence.AdminRepository
	adminSession      *persistence.AdminSessionRepository
	ticket            *persistence.TicketRepository
	ticketCategory    *persistence.TicketCategoryRepository
	ticketHistory     *persistence.TicketHistoryRepository
	maintenance       *persistence.MaintenanceRecordRepository
}

// initRepositories creates all repository instances
func initRepositories(lineClient *client.Client) *repositories {
	return &repositories{
		line:              persistence.NewLineRepository(lineClient),
		notification:      persistence.NewNotificationRepository(database.GetDB()),
		equipment:         persistence.NewEquipmentRepository(),
		brand:             persistence.NewBrandRepository(),
		equipmentCategory: persistence.NewEquipmentCategoryRepository(),
		department:        persistence.NewDepartmentRepository(),
		equipmentModel:    persistence.NewEquipmentModelRepository(),
		admin:             persistence.NewAdminRepository(),
		adminSession:      persistence.NewAdminSessionRepository(),
		ticket:            persistence.NewTicketRepository(database.GetDB()),
		ticketCategory:    persistence.NewTicketCategoryRepository(database.GetDB()),
		ticketHistory:     persistence.NewTicketHistoryRepository(database.GetDB()),
		maintenance:       persistence.NewMaintenanceRecordRepository(),
	}
}

// initUseCases creates all use case instances
func initUseCases(
	repos *repositories,
	cfg *config.Config,
	ocrClient *client.OCRClient,
	sessionStore *session.SessionStore,
) (
	*usecase.TicketUseCase,
	*usecase.MessageUseCase,
	*usecase.NotificationUseCase,
	usecase.EquipmentImportUseCase,
	usecase.AdminUsecase,
	usecase.DashboardUsecase,
	usecase.EquipmentUsecase,

) {
	equipmentMapper := mapper.NewEquipmentMapper()

	// Services
	messageService := service.NewMessageService(cfg.Contact)
	notificationService := service.NewNotificationService()
	excelParserService := service.NewExcelParserService()
	masterDataService := service.NewMasterDataService(
		repos.brand, repos.equipmentCategory, repos.department, repos.equipmentModel, equipmentMapper,
	)
	adminService := service.NewAdminService(repos.admin, repos.adminSession)
	equipmentService := service.NewEquipmentService(
		repos.equipment, repos.brand, repos.equipmentCategory, repos.department, repos.equipmentModel,
	)
	ticketNotifyService := service.NewTicketNotificationService(repos.line, repos.ticket)

	// Use cases
	ticketUC := usecase.NewTicketUseCase(
		repos.line, repos.equipment, repos.ticket, repos.ticketCategory, repos.ticketHistory, ticketNotifyService,
	)
	messageUC := usecase.NewMessageUseCase(
		repos.line, repos.equipment, ocrClient, sessionStore, messageService, ticketUC,
	)
	notificationUC := usecase.NewNotificationUseCase(repos.notification, notificationService, repos.line)
	equipmentImportUC := usecase.NewEquipmentImportUseCase(repos.equipment, excelParserService, masterDataService, equipmentMapper)
	adminUC := usecase.NewAdminUsecase(adminService)
	dashboardUC := usecase.NewDashboardUsecase(repos.equipment, repos.maintenance, repos.ticket)
	equipmentUC := usecase.NewEquipmentUsecase(equipmentService)

	return ticketUC, messageUC, notificationUC, equipmentImportUC, adminUC, dashboardUC, equipmentUC
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

	// Initialize session store for OCR confirmations
	sessionStore := session.NewSessionStore()

	// Initialize all layers
	repos := initRepositories(lineClient)
	ticketUC, messageUC, notificationUC, equipmentImportUC, adminUC, dashboardUC, equipmentUC := initUseCases(repos, cfg, ocrClient, sessionStore)

	// Initialize handlers (Interface Layer)
	webhookHandler := handlers.NewWebhookHandler(cfg.LineChannelSecret, messageUC)
	notificationHandler := handlers.NewNotificationHandler(notificationUC)
	equipmentImportHandler := handlers.NewEquipmentImportHandler(equipmentImportUC)
	adminHandler := handlers.NewAdminHandler(adminUC)
	dashboardHandler := handlers.NewDashboardHandler(dashboardUC)
	equipmentHandler := handlers.NewEquipmentHandler(equipmentUC)
	ticketHandler := handlers.NewTicketHandler(ticketUC)

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
	routes.Setup(app, webhookHandler, notificationHandler, equipmentImportHandler, adminHandler, dashboardHandler, equipmentHandler, ticketHandler)

	// Initialize and Start Notification Scheduler
	notificationScheduler := scheduler.NewNotificationScheduler(notificationUC)
	notificationScheduler.Start()
	log.Println("Notification scheduler started")

	// Cleanup function
	cleanup := func() {
		log.Println("Shutting down gracefully...")
		if notificationScheduler != nil {
			notificationScheduler.Stop()
		}
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
