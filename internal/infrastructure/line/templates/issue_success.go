package templates

import "fmt"

// GetIssueSuccessFlex returns a Flex Message confirming issue has been reported
func GetIssueSuccessFlex(serialNumber string) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#4CAF50",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "✅ บันทึกเรียบร้อย", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เลขเครื่อง: %s", serialNumber), "size": "sm", "color": "#888888",
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "text", "text": "ระบบได้รับข้อมูลการแจ้งปัญหาแล้วค่ะ", "size": "sm", "margin": "md", "wrap": true,
				},
				map[string]interface{}{
					"type": "text", "text": "เจ้าหน้าที่จะติดต่อกลับโดยเร็วที่สุด", "size": "sm", "color": "#888888", "margin": "sm", "wrap": true,
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "secondary",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "🏠 กลับหน้าหลัก",
						"data":        "action=main_menu",
						"displayText": "กลับหน้าหลัก",
					},
				},
			},
		},
	}
}
