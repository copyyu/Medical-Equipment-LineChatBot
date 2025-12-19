package bootstrap

import (
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/config"
	"medical-webhook/internal/domain/line/service"
	"medical-webhook/internal/infrastructure/database"
	"medical-webhook/internal/infrastructure/line"
	"medical-webhook/internal/interfaces/http/handlers"
	"medical-webhook/internal/interfaces/http/middleware"
	"medical-webhook/internal/interfaces/http/routes"

	"github.com/gofiber/fiber/v2"
)

type Application struct {
	Server         *fiber.App
	Config         *config.Config
	WebhookHandler *handlers.WebhookHandler
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
	lineClient, err := line.NewClient(cfg.LineChannelToken)
	if err != nil {
		return nil, nil, err
	}

	// Initialize repositories (Infrastructure Layer)
	lineRepo := line.NewRepositoryImpl(lineClient)

	// Initialize services (Domain Layer)
	messageService := service.NewMessageService()

	// Initialize use cases (Application Layer)
	messageUseCase := usecase.NewMessageUseCase(lineRepo, messageService)

	// Initialize handlers (Interface Layer)
	webhookHandler := handlers.NewWebhookHandler(cfg.LineChannelSecret, messageUseCase)

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
	routes.Setup(app, webhookHandler)

	// Cleanup function
	cleanup := func() {
		if err := database.Close(); err != nil {
			// Log error if needed
		}
	}

	return &Application{
		Server:         app,
		Config:         cfg,
		WebhookHandler: webhookHandler,
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
