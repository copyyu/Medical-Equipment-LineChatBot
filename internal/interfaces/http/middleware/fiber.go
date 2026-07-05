package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
)

// FiberMiddleware registers the global middleware chain. Order matters:
//   - RequestContext first, so a request ID exists for everything downstream.
//   - RequestLogger next, so it wraps recover and still logs panics (surfaced as
//     500s) with an accurate status.
//   - recover, to turn a panic in any handler into a 500 instead of crashing.
//   - CORS last, before the routes.
//
// allowedOrigins is the comma-separated CORS allow-list; pass a specific origin
// list in production rather than the "*" wildcard.
func FiberMiddleware(app *fiber.App, allowedOrigins string) {
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}
	app.Use(RequestContext())
	app.Use(RequestLogger())
	app.Use(recovermw.New(recovermw.Config{EnableStackTrace: true}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
}
