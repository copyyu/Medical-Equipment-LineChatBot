package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func FiberMiddleware(app *fiber.App) {
	app.Use(
		// Recover from any handler panic and turn it into a 500 via the app's
		// ErrorHandler, instead of relying on connection-level recovery (which
		// can drop the connection).
		recover.New(recover.Config{EnableStackTrace: true}),
		cors.New(cors.Config{
			AllowOrigins: "*", // need to be changed in production
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		}),
	)
}
