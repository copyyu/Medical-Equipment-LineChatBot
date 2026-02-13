package templates

import (
	"fmt"
	"medical-webhook/internal/domain/line/entity"
	"time"
)

// GetTicketStatusChangedFlex returns flex message for ticket status change notification
// This is sent as a push message to the reporter when admin changes ticket status
func GetTicketStatusChangedFlex(ticket *entity.Ticket, oldStatus, newStatus entity.TicketStatus, note string) map[string]interface{} {
	newStatusColor := newStatus.GetColor()
	newStatusText := newStatus.GetStatusText()
	oldStatusText := oldStatus.GetStatusText()
	oldStatusColor := oldStatus.GetColor()

	// Build body contents
	bodyContents := []interface{}{
		// Ticket number
		map[string]interface{}{
			"type":   "text",
			"text":   ticket.TicketNo,
			"size":   "xxl",
			"weight": "bold",
			"color":  "#333333",
			"margin": "none",
		},
		// Equipment name
		map[string]interface{}{
			"type":   "text",
			"text":   getEquipmentName(ticket),
			"size":   "md",
			"color":  "#666666",
			"wrap":   true,
			"margin": "sm",
		},
		// Separator
		map[string]interface{}{
			"type":   "separator",
			"margin": "lg",
		},
		// Status transition section
		map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"margin":  "lg",
			"spacing": "md",
			"contents": []interface{}{
				// "สถานะเปลี่ยนจาก" label
				map[string]interface{}{
					"type":   "text",
					"text":   "สถานะเปลี่ยนจาก",
					"size":   "sm",
					"color":  "#999999",
					"margin": "none",
				},
				// Old → New status row
				map[string]interface{}{
					"type":    "box",
					"layout":  "horizontal",
					"margin":  "sm",
					"spacing": "md",
					"contents": []interface{}{
						// Old status badge
						map[string]interface{}{
							"type":            "box",
							"layout":          "vertical",
							"cornerRadius":    "8px",
							"backgroundColor": oldStatusColor + "20",
							"paddingAll":      "8px",
							"flex":            0,
							"contents": []interface{}{
								map[string]interface{}{
									"type":   "text",
									"text":   oldStatusText,
									"size":   "sm",
									"color":  oldStatusColor,
									"weight": "bold",
									"align":  "center",
								},
							},
						},
						// Arrow
						map[string]interface{}{
							"type":    "box",
							"layout":  "vertical",
							"flex":    0,
							"gravity": "center",
							"contents": []interface{}{
								map[string]interface{}{
									"type":   "text",
									"text":   "→",
									"size":   "lg",
									"color":  "#999999",
									"weight": "bold",
									"align":  "center",
								},
							},
						},
						// New status badge
						map[string]interface{}{
							"type":            "box",
							"layout":          "vertical",
							"cornerRadius":    "8px",
							"backgroundColor": newStatusColor + "30",
							"paddingAll":      "8px",
							"flex":            0,
							"contents": []interface{}{
								map[string]interface{}{
									"type":   "text",
									"text":   newStatusText,
									"size":   "sm",
									"color":  newStatusColor,
									"weight": "bold",
									"align":  "center",
								},
							},
						},
					},
				},
			},
		},
		// Separator
		map[string]interface{}{
			"type":   "separator",
			"margin": "lg",
		},
		// Updated at
		createInfoRow("อัปเดตเมื่อ", time.Now().Format("2006-01-02 15:04"), "#666666"),
	}

	// Add note if provided
	if note != "" {
		bodyContents = append(bodyContents,
			map[string]interface{}{
				"type":   "separator",
				"margin": "lg",
			},
			map[string]interface{}{
				"type":    "box",
				"layout":  "horizontal",
				"margin":  "lg",
				"spacing": "sm",
				"contents": []interface{}{
					map[string]interface{}{
						"type":  "text",
						"text":  "📝",
						"size":  "sm",
						"flex":  0,
						"align": "start",
					},
					map[string]interface{}{
						"type":  "text",
						"text":  note,
						"size":  "sm",
						"color": "#666666",
						"wrap":  true,
						"flex":  1,
					},
				},
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
					"text":   fmt.Sprintf("%s อัปเดตสถานะ Ticket", getStatusEmoji(newStatus)),
					"weight": "bold",
					"size":   "lg",
					"color":  "#FFFFFF",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   newStatusText,
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

// getStatusEmoji returns emoji for status (local to this file)
func getStatusEmoji(status entity.TicketStatus) string {
	switch status {
	case entity.TicketStatusInProgress:
		return "🔧"
	case entity.TicketStatusCompleted:
		return "✅"
	case entity.TicketStatusSendToOutsource:
		return "📤"
	default:
		return "🔔"
	}
}
