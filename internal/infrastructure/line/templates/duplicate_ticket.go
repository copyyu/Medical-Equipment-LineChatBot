package templates

import "fmt"

// GetDuplicateTicketFlex returns a Flex Message notifying user about existing ticket
func GetDuplicateTicketFlex(ticketNo string, equipmentSerial string, status string) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": ColorWarning,
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "⚠️ พบรายการแจ้งซ่อมที่มีอยู่แล้ว", "color": ColorWhite, "size": "md", "weight": "bold", "wrap": true},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เลขเครื่อง: %s", equipmentSerial), "size": "sm", "color": ColorTextLight,
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "text", "text": "คุณได้แจ้งซ่อมเครื่องนี้ไปแล้ว", "size": "sm", "margin": "md", "wrap": true, "color": ColorText,
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "md",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "หมายเลข:", "size": "sm", "color": ColorTextLight, "flex": 2},
						map[string]interface{}{"type": "text", "text": ticketNo, "size": "sm", "color": ColorAccent, "weight": "bold", "flex": 3},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "สถานะ:", "size": "sm", "color": ColorTextLight, "flex": 2},
						map[string]interface{}{"type": "text", "text": status, "size": "sm", "color": ColorWarning, "weight": "bold", "flex": 3},
					},
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "primary", "color": ColorAccent,
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "📋 ดูสถานะ Ticket",
						"data":  fmt.Sprintf("action=view_ticket_status&ticket_no=%s", ticketNo),
					},
				},
				map[string]interface{}{
					"type": "button", "style": "secondary", "margin": "sm",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "🏠 กลับเมนูหลัก",
						"data":  "action=main_menu",
					},
				},
			},
		},
	}
}
