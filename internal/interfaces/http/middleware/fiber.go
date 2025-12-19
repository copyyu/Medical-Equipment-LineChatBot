package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func FiberMiddleware(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{
			AllowOrigins: "*", // need to be changed in production
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		}),
	)
}
