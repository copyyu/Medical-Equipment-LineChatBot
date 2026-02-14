package templates

import (
	"fmt"
	"medical-webhook/internal/domain/line/entity"
)

// GetTicketCreatedFlex returns flex message for ticket creation success
func GetTicketCreatedFlex(ticket *entity.Ticket) map[string]interface{} {
	statusColor := ticket.Status.GetColor()
	statusText := ticket.Status.GetStatusText()

	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#66BB6A",
			"paddingAll":      "20px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "✅ สร้าง Ticket สำเร็จ",
					"weight": "bold",
					"size":   "xl",
					"color":  "#FFFFFF",
				},
			},
		},
		"body": map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"spacing": "md",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "หมายเลข Ticket",
					"size":   "sm",
					"color":  "#999999",
					"margin": "none",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   ticket.TicketNo,
					"size":   "xxl",
					"weight": "bold",
					"color":  "#1DB446",
					"margin": "sm",
				},
				map[string]interface{}{
					"type":   "separator",
					"margin": "lg",
				},
				map[string]interface{}{
					"type":    "box",
					"layout":  "vertical",
					"margin":  "lg",
					"spacing": "sm",
					"contents": []interface{}{
						createInfoRow("สถานะ", statusText, statusColor),
						createInfoRow("อุปกรณ์", getEquipmentName(ticket), "#333333"),
						createInfoRow("วันที่แจ้ง", ticket.ReportedAt.Format("2006-01-02 15:04"), "#666666"),
					},
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "📋 บันทึกเลข Ticket นี้ไว้เพื่อติดตามสถานะ\nหรือใช้เมนู 'ติดตามสถานะ'",
					"size":   "xs",
					"color":  "#999999",
					"margin": "lg",
					"wrap":   true,
				},
			},
		},
		"footer": map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "button",
					"style": "primary",
					"color": "#42A5F5",
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
				map[string]interface{}{
					"type":  "button",
					"style": "link",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "กลับเมนูหลัก",
						"data":  "action=main_menu",
					},
				},
			},
		},
	}
}

// GetTicketStatusFlex returns flex message for ticket status inquiry
func GetTicketStatusFlex(ticket *entity.Ticket) map[string]interface{} {
	statusColor := ticket.Status.GetColor()
	statusText := ticket.Status.GetStatusText()
	priorityColor := ticket.Priority.GetColor()
	priorityText := ticket.Priority.GetPriorityText()

	bodyContents := []interface{}{
		map[string]interface{}{
			"type":   "text",
			"text":   getEquipmentName(ticket),
			"size":   "lg",
			"weight": "bold",
			"wrap":   true,
			"margin": "none",
		},
	}

	// Add description if exists
	if ticket.Description != nil && *ticket.Description != "" {
		bodyContents = append(bodyContents, map[string]interface{}{
			"type":   "text",
			"text":   *ticket.Description,
			"size":   "sm",
			"color":  "#666666",
			"wrap":   true,
			"margin": "md",
		})
	}

	bodyContents = append(bodyContents,
		map[string]interface{}{
			"type":   "separator",
			"margin": "md",
		},
		map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"margin":  "lg",
			"spacing": "sm",
			"contents": []interface{}{
				createInfoRow("สถานะ", statusText, statusColor),
				createInfoRow("ความสำคัญ", priorityText, priorityColor),
				createInfoRow("วันที่แจ้ง", ticket.ReportedAt.Format("2006-01-02 15:04"), "#666666"),
			},
		},
	)

	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": statusColor,
			"paddingAll":      "20px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "📋 สถานะ Ticket",
					"weight": "bold",
					"size":   "xl",
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
					"color": "#42A5F5",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "รายการของฉัน",
						"data":  "action=my_tickets",
					},
				},
				map[string]interface{}{
					"type":  "button",
					"style": "link",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "กลับเมนูหลัก",
						"data":  "action=main_menu",
					},
				},
			},
		},
	}
}

// GetMyTicketsFlex returns flex message showing user's tickets as carousel
func GetMyTicketsFlex(tickets []entity.Ticket) map[string]interface{} {
	bubbles := []interface{}{}

	// Show up to 10 tickets
	maxTickets := len(tickets)
	if maxTickets > 10 {
		maxTickets = 10
	}

	for i := 0; i < maxTickets; i++ {
		ticket := tickets[i]
		jobStatusColor := ticket.Status.GetColor()
		jobStatusText := ticket.Status.GetStatusText()
		priorityText := ticket.Priority.GetPriorityText()
		priorityColor := ticket.Priority.GetColor()

		// Get asset status from equipment relation
		assetStatusText := "ไม่ทราบ"
		assetStatusColor := "#78909C"
		equipmentCode := ""
		if ticket.Equipment != nil {
			assetStatusText = ticket.Equipment.Status.GetStatusText()
			assetStatusColor = ticket.Equipment.Status.GetColor()
			equipmentCode = ticket.Equipment.IDCode
		}

		// Category name
		categoryName := "ไม่ระบุ"
		if ticket.Category.Name != "" {
			categoryName = ticket.Category.Name
		}

		// Equipment display name
		equipName := getEquipmentName(&ticket)

		// Build body contents
		bodyContents := []interface{}{
			// Equipment name
			map[string]interface{}{
				"type":   "text",
				"text":   equipName,
				"weight": "bold",
				"size":   "md",
				"wrap":   true,
				"color":  "#333333",
			},
		}

		// Equipment code (if available)
		if equipmentCode != "" {
			bodyContents = append(bodyContents, map[string]interface{}{
				"type":   "text",
				"text":   fmt.Sprintf("🏷️ %s", equipmentCode),
				"size":   "xs",
				"color":  "#0367D3",
				"margin": "sm",
			})
		}

		// Separator
		bodyContents = append(bodyContents, map[string]interface{}{
			"type":   "separator",
			"margin": "lg",
		})

		// Status badges section
		bodyContents = append(bodyContents,
			// Job Status row
			map[string]interface{}{
				"type":    "box",
				"layout":  "horizontal",
				"margin":  "lg",
				"spacing": "sm",
				"contents": []interface{}{
					map[string]interface{}{
						"type":  "text",
						"text":  "📋 สถานะงาน",
						"size":  "xs",
						"color": "#888888",
						"flex":  4,
					},
					map[string]interface{}{
						"type":            "box",
						"layout":          "vertical",
						"flex":            5,
						"cornerRadius":    "4px",
						"backgroundColor": jobStatusColor,
						"paddingAll":      "4px",
						"contents": []interface{}{
							map[string]interface{}{
								"type":  "text",
								"text":  jobStatusText,
								"size":  "xs",
								"color": "#FFFFFF",
								"align": "center",
							},
						},
					},
				},
			},
			// Asset Status row
			map[string]interface{}{
				"type":    "box",
				"layout":  "horizontal",
				"margin":  "sm",
				"spacing": "sm",
				"contents": []interface{}{
					map[string]interface{}{
						"type":  "text",
						"text":  "🔧 สถานะเครื่อง",
						"size":  "xs",
						"color": "#888888",
						"flex":  4,
					},
					map[string]interface{}{
						"type":            "box",
						"layout":          "vertical",
						"flex":            5,
						"cornerRadius":    "4px",
						"backgroundColor": assetStatusColor,
						"paddingAll":      "4px",
						"contents": []interface{}{
							map[string]interface{}{
								"type":  "text",
								"text":  assetStatusText,
								"size":  "xs",
								"color": "#FFFFFF",
								"align": "center",
							},
						},
					},
				},
			},
		)

		// Separator
		bodyContents = append(bodyContents, map[string]interface{}{
			"type":   "separator",
			"margin": "lg",
		})

		// Info rows: category, priority, date
		bodyContents = append(bodyContents,
			// Category
			map[string]interface{}{
				"type":    "box",
				"layout":  "horizontal",
				"margin":  "lg",
				"spacing": "sm",
				"contents": []interface{}{
					map[string]interface{}{
						"type":  "text",
						"text":  "หมวดหมู่",
						"size":  "xs",
						"color": "#888888",
						"flex":  4,
					},
					map[string]interface{}{
						"type":   "text",
						"text":   categoryName,
						"size":   "xs",
						"color":  "#333333",
						"flex":   5,
						"weight": "bold",
					},
				},
			},
			// Priority
			map[string]interface{}{
				"type":    "box",
				"layout":  "horizontal",
				"margin":  "sm",
				"spacing": "sm",
				"contents": []interface{}{
					map[string]interface{}{
						"type":  "text",
						"text":  "ความสำคัญ",
						"size":  "xs",
						"color": "#888888",
						"flex":  4,
					},
					map[string]interface{}{
						"type":   "text",
						"text":   fmt.Sprintf("● %s", priorityText),
						"size":   "xs",
						"color":  priorityColor,
						"flex":   5,
						"weight": "bold",
					},
				},
			},
			// Date
			map[string]interface{}{
				"type":    "box",
				"layout":  "horizontal",
				"margin":  "sm",
				"spacing": "sm",
				"contents": []interface{}{
					map[string]interface{}{
						"type":  "text",
						"text":  "วันที่แจ้ง",
						"size":  "xs",
						"color": "#888888",
						"flex":  4,
					},
					map[string]interface{}{
						"type":  "text",
						"text":  ticket.ReportedAt.Format("02/01/2006 15:04"),
						"size":  "xs",
						"color": "#333333",
						"flex":  5,
					},
				},
			},
		)

		bubble := map[string]interface{}{
			"type": "bubble",
			"size": "kilo",
			"header": map[string]interface{}{
				"type":            "box",
				"layout":          "horizontal",
				"backgroundColor": jobStatusColor,
				"paddingAll":      "15px",
				"contents": []interface{}{
					map[string]interface{}{
						"type":   "text",
						"text":   fmt.Sprintf("🔧 %s", ticket.TicketNo),
						"weight": "bold",
						"size":   "sm",
						"color":  "#FFFFFF",
						"flex":   6,
					},
					map[string]interface{}{
						"type":            "box",
						"layout":          "vertical",
						"flex":            4,
						"cornerRadius":    "12px",
						"backgroundColor": "#FFFFFF33",
						"paddingAll":      "4px",
						"contents": []interface{}{
							map[string]interface{}{
								"type":  "text",
								"text":  jobStatusText,
								"size":  "xxs",
								"color": "#FFFFFF",
								"align": "center",
							},
						},
					},
				},
			},
			"body": map[string]interface{}{
				"type":       "box",
				"layout":     "vertical",
				"paddingAll": "15px",
				"contents":   bodyContents,
			},
			"footer": map[string]interface{}{
				"type":       "box",
				"layout":     "vertical",
				"paddingAll": "10px",
				"spacing":    "sm",
				"contents": []interface{}{
					map[string]interface{}{
						"type":   "button",
						"style":  "primary",
						"color":  jobStatusColor,
						"height": "sm",
						"action": map[string]interface{}{
							"type":  "postback",
							"label": "📄 ดูรายละเอียด",
							"data":  fmt.Sprintf("action=view_ticket&ticket_no=%s", ticket.TicketNo),
						},
					},
				},
			},
		}

		bubbles = append(bubbles, bubble)
	}

	return map[string]interface{}{
		"type":     "carousel",
		"contents": bubbles,
	}
}

// Helper function to create info row
func createInfoRow(label, value, valueColor string) map[string]interface{} {
	return map[string]interface{}{
		"type":    "box",
		"layout":  "horizontal",
		"margin":  "md",
		"spacing": "sm",
		"contents": []interface{}{
			map[string]interface{}{
				"type":  "text",
				"text":  label,
				"size":  "sm",
				"color": "#999999",
				"flex":  0,
			},
			map[string]interface{}{
				"type":   "text",
				"text":   value,
				"size":   "sm",
				"color":  valueColor,
				"weight": "bold",
				"flex":   0,
				"wrap":   true,
			},
		},
	}
}

// Helper function to get equipment name safely
func getEquipmentName(ticket *entity.Ticket) string {
	if ticket.EquipmentName != nil {
		return *ticket.EquipmentName
	}
	return "ไม่ระบุอุปกรณ์"
}
