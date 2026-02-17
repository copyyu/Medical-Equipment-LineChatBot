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

// GetMyTicketsFlex returns flex message showing user's tickets as vertical list in single bubble
func GetMyTicketsFlex(tickets []entity.Ticket) map[string]interface{} {
	// Show up to 10 tickets
	maxTickets := len(tickets)
	if maxTickets > 10 {
		maxTickets = 10
	}

	bodyContents := []interface{}{
		map[string]interface{}{
			"type": "text", "text": fmt.Sprintf("📊 รวม %d รายการ", maxTickets),
			"size": "xs", "color": "#888888", "margin": "sm",
		},
	}

	for i := 0; i < maxTickets; i++ {
		ticket := tickets[i]
		jobStatusColor := ticket.Status.GetColor()
		jobStatusText := ticket.Status.GetStatusText()
		equipName := getEquipmentName(&ticket)

		bgColor := "#FFFFFF"
		if i%2 == 1 {
			bgColor = "#F5F5F5"
		}

		bodyContents = append(bodyContents, map[string]interface{}{
			"type": "box", "layout": "vertical", "margin": "sm",
			"paddingAll": "10px", "backgroundColor": bgColor, "cornerRadius": "6px",
			"contents": []interface{}{
				// Row 1: running number + ticket_no + status badge
				map[string]interface{}{
					"type": "box", "layout": "horizontal",
					"contents": []interface{}{
						map[string]interface{}{
							"type": "text", "text": fmt.Sprintf("%d. %s", i+1, ticket.TicketNo),
							"size": "xs", "weight": "bold", "flex": 5,
						},
						map[string]interface{}{
							"type": "box", "layout": "vertical", "flex": 4,
							"cornerRadius": "4px", "backgroundColor": jobStatusColor, "paddingAll": "3px",
							"contents": []interface{}{
								map[string]interface{}{
									"type": "text", "text": jobStatusText,
									"size": "xxs", "color": "#FFFFFF", "align": "center",
								},
							},
						},
					},
				},
				// Row 2: equipment name
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("   🔧 %s", equipName),
					"size": "xxs", "color": "#666666", "wrap": true,
				},
				// Row 3: date
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("   📅 %s", ticket.ReportedAt.Format("02/01/2006 15:04")),
					"size": "xxs", "color": "#888888",
				},
			},
		})
	}

	return map[string]interface{}{
		"type": "bubble", "size": "mega",
		"header": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"backgroundColor": "#1565C0", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "รายการแจ้งปัญหาของฉัน",
					"color": "#FFFFFF", "size": "lg", "weight": "bold",
				},
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("แสดง %d รายการล่าสุด", maxTickets),
					"color": "#FFFFFFCC", "size": "xs",
				},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"spacing": "sm", "paddingAll": "12px",
			"contents": bodyContents,
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "link",
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
