package handlers

import (
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{
		adminUsecase: adminUsecase,
	}
}

func (h *AdminHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest(c, "Invalid request data")
	}

	admin, err := h.adminUsecase.Register(c.Context(), &req)
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Created(c, admin, "Admin registered successfully")
}

func (h *AdminHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest(c, "Invalid request data")
	}

	ipAddress := c.IP()
	loginResponse, err := h.adminUsecase.Login(c.Context(), &req, ipAddress)
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, loginResponse, "Login successful")
}

func (h *AdminHandler) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return errors.Unauthorized(c, "No token provided")
	}

	// Remove "Bearer " prefix
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.adminUsecase.Logout(c.Context(), token); err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, nil, "Logout successful")
}

func (h *AdminHandler) GetProfile(c *fiber.Ctx) error {
	adminID := c.Locals("admin_id").(string)

	profile, err := h.adminUsecase.GetProfile(c.Context(), adminID)
	if err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, profile, "Profile retrieved successfully")
}

func (h *AdminHandler) UpdateProfile(c *fiber.Ctx) error {
	adminID := c.Locals("admin_id").(string)

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest(c, "Invalid request data")
	}

	if err := h.adminUsecase.UpdateProfile(c.Context(), adminID, &req); err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, nil, "Profile updated successfully")
}

func (h *AdminHandler) ChangePassword(c *fiber.Ctx) error {
	adminID := c.Locals("admin_id").(string)

	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest(c, "Invalid request data")
	}

	if err := h.adminUsecase.ChangePassword(c.Context(), adminID, &req); err != nil {
		return errors.Error(c, err)
	}

	return errors.Success(c, nil, "Password changed successfully")
}
