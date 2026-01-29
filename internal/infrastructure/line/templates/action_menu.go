package templates

import "fmt"

// GetActionMenuFlex returns a Flex Message with options to view equipment info or report issue
// Shown after OCR confirmation
func GetActionMenuFlex(serialNumber string) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#4CAF50",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "✅ ยืนยันเลขเครื่องสำเร็จ", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เลขเครื่อง: %s", serialNumber), "size": "md", "weight": "bold", "color": "#0367D3",
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "text", "text": "ต้องการทำอะไรต่อ?", "size": "sm", "margin": "md",
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "primary", "color": "#5B9BD5",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "🔍 ดูข้อมูลเครื่อง",
						"data":        fmt.Sprintf("action=view_equipment_info&serial=%s", serialNumber),
						"displayText": "ดูข้อมูลเครื่อง",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "primary", "color": "#FF5722", "margin": "sm",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "⚠️ แจ้งปัญหา",
						"data":        fmt.Sprintf("action=start_report_issue&serial=%s", serialNumber),
						"displayText": "แจ้งปัญหา",
					},
				},
			},
		},
	}
}
