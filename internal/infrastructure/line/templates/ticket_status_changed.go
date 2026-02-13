package templates

import (
	"fmt"
	"medical-webhook/internal/domain/line/entity"
	"time"
)

// GetTicketStatusChangedFlex returns flex message for ticket status change notification
func GetTicketStatusChangedFlex(ticket *entity.Ticket, oldStatus, newStatus entity.TicketStatus, note string) map[string]interface{} {
	newStatusColor := newStatus.GetColor()
	newStatusText := newStatus.GetStatusText()
	oldStatusText := oldStatus.GetStatusText()

	// Build body contents
	bodyContents := []interface{}{
		// Equipment name
		map[string]interface{}{
			"type":   "text",
			"text":   getEquipmentName(ticket),
			"size":   "lg",
			"weight": "bold",
			"wrap":   true,
			"margin": "none",
		},
		// Separator
		map[string]interface{}{
			"type":   "separator",
			"margin": "lg",
		},
		// Status transition info rows
		map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"margin":  "lg",
			"spacing": "sm",
			"contents": []interface{}{
				createInfoRow("สถานะเดิม", oldStatusText, "#999999"),
				createInfoRow("สถานะใหม่", newStatusText, newStatusColor),
				createInfoRow("อัปเดตเมื่อ", time.Now().Format("2006-01-02 15:04"), "#666666"),
			},
		},
	}

	// Add note if provided
	if note != "" {
		bodyContents = append(bodyContents,
			map[string]interface{}{
				"type":   "separator",
				"margin": "lg",
			},
			map[string]interface{}{
				"type":   "text",
				"text":   fmt.Sprintf("📝 %s", note),
				"size":   "sm",
				"color":  "#666666",
				"wrap":   true,
				"margin": "lg",
			},
		)
	}

	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": newStatusColor,
			"paddingAll":      "20px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   fmt.Sprintf("%s อัปเดตสถานะ Ticket", getStatusChangeEmoji(newStatus)),
					"weight": "bold",
					"size":   "lg",
					"color":  "#FFFFFF",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   ticket.TicketNo,
					"size":   "sm",
					"color":  "#FFFFFF",
					"margin": "sm",
				},
			},
		},
		"body": map[string]interface{}{
			"type":     "box",
			"layout":   "vertical",
			"spacing":  "md",
			"contents": bodyContents,
		},
		"footer": map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "button",
					"style": "primary",
					"color": newStatusColor,
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "ดูรายละเอียด",
						"data":  fmt.Sprintf("action=view_ticket&ticket_no=%s", ticket.TicketNo),
					},
				},
				map[string]interface{}{
					"type":  "button",
					"style": "link",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "รายการของฉัน",
						"data":  "action=my_tickets",
					},
				},
			},
		},
	}
}

// getStatusChangeEmoji returns emoji for status
func getStatusChangeEmoji(status entity.TicketStatus) string {
	switch status {
	case entity.TicketStatusInProcess:
		return "🔧"
	case entity.TicketStatusCompleted:
		return "✅"
	case entity.TicketStatusSendToOutsource:
		return "📤"
	default:
		return "🔔"
	}
}
