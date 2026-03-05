package handlers

import (
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/usecase"
	"strconv"

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

// DownloadExpiryExcel - สร้างและดาวน์โหลดไฟล์ Excel รายการเครื่องใกล้หมดอายุ
func (h *NotificationHandler) DownloadExpiryExcel(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse optional dept_id
	var departmentID *uint
	if deptIDStr := c.Query("dept_id"); deptIDStr != "" {
		id, err := strconv.ParseUint(deptIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid dept_id parameter",
			})
		}
		deptID := uint(id)
		departmentID = &deptID
	}
	filter := c.Query("filter", "all") // this_year, next_year, all

	xlsxBytes, filename, err := h.notificationUseCase.BuildExpiryExcel(ctx, departmentID, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	return c.Send(xlsxBytes)
}
