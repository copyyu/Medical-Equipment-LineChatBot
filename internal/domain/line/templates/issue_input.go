package templates

import "fmt"

// GetIssueInputFlex returns a Flex Message asking whether to input issue description or skip
func GetIssueInputFlex(serialNumber string) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#FF5722",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "⚠️ แจ้งปัญหา", "color": "#FFFFFF", "size": "md", "weight": "bold"},
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
					"type": "text", "text": "ต้องการระบุรายละเอียดปัญหาหรือไม่?", "size": "sm", "margin": "md", "wrap": true,
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "primary", "color": "#4CAF50",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "✏️ พิมพ์รายละเอียด",
						"data":        fmt.Sprintf("action=input_issue_desc&serial=%s", serialNumber),
						"displayText": "พิมพ์รายละเอียดปัญหา",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "secondary", "margin": "sm",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "⏭️ ข้าม (ไม่ระบุ)",
						"data":        fmt.Sprintf("action=submit_issue&serial=%s&desc=", serialNumber),
						"displayText": "ข้ามการระบุรายละเอียด",
					},
				},
			},
		},
	}
}
