package bootstrap

import (
	"log"
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/config"
	"medical-webhook/internal/domain/line/service"
	"medical-webhook/internal/infrastructure/client"
	"medical-webhook/internal/infrastructure/database"
	"medical-webhook/internal/infrastructure/line"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"
	"medical-webhook/internal/interfaces/http/routes"
	"medical-webhook/internal/interfaces/http/scheduler"

	"github.com/gofiber/fiber/v2"
)

type Application struct {
	Server                *fiber.App
	Config                *config.Config
	WebhookHandler        *handlers.WebhookHandler
	NotificationHandler   *handlers.NotificationHandler
	NotificationScheduler *scheduler.NotificationScheduler
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

	// Initialize repositories (Infrastructure Layer)
	lineRepo := line.NewLineRepository(lineClient)
	notificationRepo := line.NewNotificationRepository(database.GetDB())

	// Initialize services (Domain Layer)
	messageService := service.NewMessageService()
	notificationService := service.NewNotificationService()

	// Initialize use cases (Application Layer)
	messageUseCase := usecase.NewMessageUseCase(lineRepo, messageService)
	notificationUseCase := usecase.NewNotificationUseCase(
		notificationRepo,
		notificationService,
		lineRepo,
	)

	// Initialize handlers (Interface Layer)
	webhookHandler := handlers.NewWebhookHandler(cfg.LineChannelSecret, messageUseCase)
	notificationHandler := handlers.NewNotificationHandler(notificationUseCase)

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
	routes.Setup(app, webhookHandler, notificationHandler)

	// Initialize และ Start Notification Scheduler
	notificationScheduler := scheduler.NewNotificationScheduler(notificationUseCase)
	notificationScheduler.Start()
	log.Println("Notification scheduler started")

	// Cleanup function
	cleanup := func() {
		log.Println("Shutting down gracefully...")
		//  Stop scheduler
		if notificationScheduler != nil {
			notificationScheduler.Stop()
		}
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		log.Println("Cleanup complete")
	}

	return &Application{
		Server:              app,
		Config:              cfg,
		WebhookHandler:      webhookHandler,
		NotificationHandler: notificationHandler,
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
