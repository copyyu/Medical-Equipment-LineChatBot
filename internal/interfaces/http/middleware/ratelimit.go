package middleware

import (
	"time"

	"medical-webhook/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// AuthRateLimiter throttles authentication endpoints (login/register) to blunt
// brute-force and abuse. Default: 10 requests per minute per client IP. The
// limit-reached response uses the standard error envelope (429 RATE_LIMITED).
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return errors.Error(c, fiber.NewError(fiber.StatusTooManyRequests, "Too many requests, please slow down"))
		},
	})
}
