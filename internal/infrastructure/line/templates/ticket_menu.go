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
			"backgroundColor": ColorSuccess,
			"paddingAll":      "20px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "✅ สร้าง Ticket สำเร็จ",
					"weight": "bold",
					"size":   "xl",
					"color":  ColorWhite,
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
					"color":  ColorTextLight,
					"margin": "none",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   ticket.TicketNo,
					"size":   "xxl",
					"weight": "bold",
					"color":  ColorSuccess,
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
						createInfoRow("อุปกรณ์", getEquipmentName(ticket), ColorText),
						createInfoRow("วันที่แจ้ง", ticket.ReportedAt.Format("2006-01-02 15:04"), ColorTextLight),
					},
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "📋 บันทึกเลข Ticket นี้ไว้เพื่อติดตามสถานะ\nหรือใช้เมนู 'ติดตามสถานะ'",
					"size":   "xs",
					"color":  ColorTextLight,
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
					"color": ColorAccent,
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
			"color":  ColorText,
		},
	}

	// Add description if exists
	if ticket.Description != nil && *ticket.Description != "" {
		bodyContents = append(bodyContents, map[string]interface{}{
			"type":   "text",
			"text":   *ticket.Description,
			"size":   "sm",
			"color":  ColorTextLight,
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
				createInfoRow("วันที่แจ้ง", ticket.ReportedAt.Format("2006-01-02 15:04"), ColorTextLight),
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
					"color":  ColorWhite,
				},
				map[string]interface{}{
					"type":   "text",
					"text":   ticket.TicketNo,
					"size":   "sm",
					"color":  ColorWhite,
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
					"color": ColorAccent,
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
			"size": "xs", "color": ColorTextLight, "margin": "sm",
		},
	}

	for i := 0; i < maxTickets; i++ {
		ticket := tickets[i]
		jobStatusColor := ticket.Status.GetColor()
		jobStatusText := ticket.Status.GetStatusText()
		equipName := getEquipmentName(&ticket)

		bgColor := ColorWhite
		if i%2 == 1 {
			bgColor = ColorBgAlt
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
							"size": "xs", "weight": "bold", "flex": 5, "color": ColorText,
						},
						map[string]interface{}{
							"type": "box", "layout": "vertical", "flex": 4,
							"cornerRadius": "4px", "backgroundColor": jobStatusColor, "paddingAll": "3px",
							"contents": []interface{}{
								map[string]interface{}{
									"type": "text", "text": jobStatusText,
									"size": "xxs", "color": ColorWhite, "align": "center",
								},
							},
						},
					},
				},
				// Row 2: equipment name
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("   🔧 %s", equipName),
					"size": "xxs", "color": ColorTextLight, "wrap": true,
				},
				// Row 3: date
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("   📅 %s", ticket.ReportedAt.Format("02/01/2006 15:04")),
					"size": "xxs", "color": ColorTextLight,
				},
			},
		})
	}

	return map[string]interface{}{
		"type": "bubble", "size": "mega",
		"header": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"backgroundColor": ColorPrimaryDark, "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "📋 รายการแจ้งปัญหาของฉัน",
					"color": ColorWhite, "size": "lg", "weight": "bold",
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
				"color": ColorTextLight,
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

// GetTicketStatusFilterFlex returns flex message for selecting ticket status filter
func GetTicketStatusFilterFlex() map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": ColorPrimaryDark,
			"paddingAll":      "20px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "🔍 เลือกสถานะ",
					"weight": "bold",
					"size":   "xl",
					"color":  ColorWhite,
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "กรุณาเลือกสถานะงานที่ต้องการดู",
					"size":   "sm",
					"color":  "#FFFFFFCC",
					"margin": "sm",
				},
			},
		},
		"body": map[string]interface{}{
			"type":       "box",
			"layout":     "vertical",
			"spacing":    "md",
			"paddingAll": "20px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "button",
					"style": "primary",
					"color": entity.TicketStatusInProcess.GetColor(),
					"action": map[string]interface{}{
						"type":  "postback",
						"label": entity.TicketStatusInProcess.GetStatusText(),
						"data":  "action=filter_tickets&status=" + string(entity.TicketStatusInProcess),
					},
				},
				map[string]interface{}{
					"type":  "button",
					"style": "primary",
					"color": entity.TicketStatusSendToOutsource.GetColor(),
					"action": map[string]interface{}{
						"type":  "postback",
						"label": entity.TicketStatusSendToOutsource.GetStatusText(),
						"data":  "action=filter_tickets&status=" + string(entity.TicketStatusSendToOutsource),
					},
				},
				map[string]interface{}{
					"type":  "button",
					"style": "primary",
					"color": entity.TicketStatusCompleted.GetColor(),
					"action": map[string]interface{}{
						"type":  "postback",
						"label": entity.TicketStatusCompleted.GetStatusText(),
						"data":  "action=filter_tickets&status=" + string(entity.TicketStatusCompleted),
					},
				},
				map[string]interface{}{
					"type":  "button",
					"style": "secondary",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "ดูทั้งหมด",
						"data":  "action=filter_tickets&status=ALL",
					},
				},
			},
		},
	}
}
