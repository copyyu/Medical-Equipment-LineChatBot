package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success response
func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created response
func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error response with auto status code mapping
func Error(c *fiber.Ctx, err error) error {
	statusCode, message := MapErrorToResponse(err)
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// MapErrorToResponse maps domain errors to HTTP status codes
func MapErrorToResponse(err error) (int, string) {
	switch {
	// Admin errors
	case errors.Is(err, ErrInvalidCredentials):
		return fiber.StatusUnauthorized, "Invalid username or password"
	case errors.Is(err, ErrAdminInactive):
		return fiber.StatusForbidden, "Account is inactive"
	case errors.Is(err, ErrInvalidToken):
		return fiber.StatusUnauthorized, "Invalid or expired token"
	case errors.Is(err, ErrUsernameExists):
		return fiber.StatusConflict, "Username already exists"
	case errors.Is(err, ErrEmailExists):
		return fiber.StatusConflict, "Email already exists"
	case errors.Is(err, ErrAdminNotFound):
		return fiber.StatusNotFound, "Admin not found"
	case errors.Is(err, ErrSessionExpired):
		return fiber.StatusUnauthorized, "Session expired"
	case errors.Is(err, ErrWeakPassword):
		return fiber.StatusBadRequest, "Password is too weak"

	// // Common errors
	case errors.Is(err, ErrNotFound):
		return fiber.StatusNotFound, "Resource not found"
	case errors.Is(err, ErrUnauthorized):
		return fiber.StatusUnauthorized, "Unauthorized"
	case errors.Is(err, ErrForbidden):
		return fiber.StatusForbidden, "Forbidden"
	case errors.Is(err, ErrBadRequest):
		return fiber.StatusBadRequest, "Bad request"
	case errors.Is(err, ErrConflict):
		return fiber.StatusConflict, "Resource conflict"
	case errors.Is(err, ErrValidationFailed):
		return fiber.StatusBadRequest, "Validation failed"

	// Default
	default:
		return fiber.StatusInternalServerError, "Internal server error"
	}
}

// BadRequest response
func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// Unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// NotFound response
func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// InternalServerError response
func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Error:   message,
	})
}
