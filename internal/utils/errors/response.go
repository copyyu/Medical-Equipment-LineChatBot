package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Code      string      `json:"code,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// requestID returns the current request's ID (set by the RequestContext
// middleware into the X-Request-ID response header) for client/log correlation.
func requestID(c *fiber.Ctx) string {
	return c.GetRespHeader(fiber.HeaderXRequestID)
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

// Error writes the standard error envelope with an auto-mapped HTTP status, a
// stable machine-readable code, and the request ID for correlation.
func Error(c *fiber.Ctx, err error) error {
	statusCode, code, message := MapErrorToResponse(err)
	return c.Status(statusCode).JSON(Response{
		Success:   false,
		Error:     message,
		Code:      code,
		RequestID: requestID(c),
	})
}

// FiberErrorHandler is the app-level error handler: it renders any error that
// reaches Fiber uncaught (router 404s, panics surfaced by recover, handler
// errors) using the same standard envelope, without leaking internal details.
func FiberErrorHandler(c *fiber.Ctx, err error) error {
	return Error(c, err)
}

// MapErrorToResponse maps an error to (httpStatus, machineCode, humanMessage).
func MapErrorToResponse(err error) (int, string, string) {
	// Respect explicit Fiber errors (fiber.NewError, router 404s, etc.).
	var fe *fiber.Error
	if errors.As(err, &fe) {
		return fe.Code, codeForStatus(fe.Code), fe.Message
	}

	switch {
	// Admin errors
	case errors.Is(err, ErrInvalidCredentials):
		return fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password"
	case errors.Is(err, ErrAdminInactive):
		return fiber.StatusForbidden, "ACCOUNT_INACTIVE", "Account is inactive"
	case errors.Is(err, ErrInvalidToken):
		return fiber.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token"
	case errors.Is(err, ErrUsernameExists):
		return fiber.StatusConflict, "USERNAME_EXISTS", "Username already exists"
	case errors.Is(err, ErrEmailExists):
		return fiber.StatusConflict, "EMAIL_EXISTS", "Email already exists"
	case errors.Is(err, ErrAdminNotFound):
		return fiber.StatusNotFound, "ADMIN_NOT_FOUND", "Admin not found"
	case errors.Is(err, ErrSessionExpired):
		return fiber.StatusUnauthorized, "SESSION_EXPIRED", "Session expired"
	case errors.Is(err, ErrWeakPassword):
		return fiber.StatusBadRequest, "WEAK_PASSWORD", "Password is too weak"

	// Common errors
	case errors.Is(err, ErrNotFound):
		return fiber.StatusNotFound, "NOT_FOUND", "Resource not found"
	case errors.Is(err, ErrUnauthorized):
		return fiber.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"
	case errors.Is(err, ErrForbidden):
		return fiber.StatusForbidden, "FORBIDDEN", "Forbidden"
	case errors.Is(err, ErrBadRequest):
		return fiber.StatusBadRequest, "BAD_REQUEST", "Bad request"
	case errors.Is(err, ErrConflict):
		return fiber.StatusConflict, "CONFLICT", "Resource conflict"
	case errors.Is(err, ErrValidationFailed):
		return fiber.StatusBadRequest, "VALIDATION_FAILED", "Validation failed"

	// Default: do not leak internal error details.
	default:
		return fiber.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"
	}
}

// codeForStatus returns a stable machine code for a bare HTTP status.
func codeForStatus(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return "BAD_REQUEST"
	case fiber.StatusUnauthorized:
		return "UNAUTHORIZED"
	case fiber.StatusForbidden:
		return "FORBIDDEN"
	case fiber.StatusNotFound:
		return "NOT_FOUND"
	case fiber.StatusConflict:
		return "CONFLICT"
	case fiber.StatusTooManyRequests:
		return "RATE_LIMITED"
	default:
		if status >= 500 {
			return "INTERNAL_ERROR"
		}
		return "ERROR"
	}
}

// BadRequest response
func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success:   false,
		Error:     message,
		Code:      "BAD_REQUEST",
		RequestID: requestID(c),
	})
}

// Unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success:   false,
		Error:     message,
		Code:      "UNAUTHORIZED",
		RequestID: requestID(c),
	})
}

// Forbidden response
func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success:   false,
		Error:     message,
		Code:      "FORBIDDEN",
		RequestID: requestID(c),
	})
}

// NotFound response
func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success:   false,
		Error:     message,
		Code:      "NOT_FOUND",
		RequestID: requestID(c),
	})
}

// InternalServerError response
func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success:   false,
		Error:     message,
		Code:      "INTERNAL_ERROR",
		RequestID: requestID(c),
	})
}
