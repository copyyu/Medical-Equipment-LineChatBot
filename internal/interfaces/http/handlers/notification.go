package handlers

import (
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/usecase"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notificationUseCase *usecase.NotificationUseCase
}

func NewNotificationHandler(notificationUseCase *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: notificationUseCase,
	}
}

// SendJuneAlerts - ส่งการแจ้งเตือนรอบมิถุนายน (Manual)
func (h *NotificationHandler) SendJuneAlerts(c *fiber.Ctx) error {
	err := h.notificationUseCase.SendJuneAlerts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "June alerts sent successfully",
	})
}

// SendAugustAlerts - ส่งการแจ้งเตือนรอบสิงหาคม (Manual)
func (h *NotificationHandler) SendAugustAlerts(c *fiber.Ctx) error {
	err := h.notificationUseCase.SendAugustAlerts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "August alerts sent successfully",
	})
}

// GetSummary - สรุปการแจ้งเตือน
func (h *NotificationHandler) GetSummary(c *fiber.Ctx) error {
	summary, err := h.notificationUseCase.GetNotificationSummary(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}

// UpdateSettings - อัพเดทการตั้งค่า
func (h *NotificationHandler) UpdateSettings(c *fiber.Ctx) error {
	var settingsDTO dto.NotificationSettingDTO
	if err := c.BodyParser(&settingsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.notificationUseCase.UpdateSettings(c.Context(), &settingsDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Settings updated successfully",
	})
}
