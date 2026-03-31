package templates

import "fmt"

// GetIssueInputFlex returns a Flex Message asking whether to input issue description or skip
func GetIssueInputFlex(serialNumber string, categoryID uint) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": ColorDanger,
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "⚠️ แจ้งปัญหา", "color": ColorWhite, "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เลขเครื่อง: %s", serialNumber), "size": "sm", "color": ColorTextLight,
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "text", "text": "ต้องการระบุรายละเอียดปัญหาหรือไม่?", "size": "sm", "margin": "md", "wrap": true, "color": ColorText,
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "primary", "color": ColorSuccess,
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "✏️ พิมพ์รายละเอียด",
						"data":        fmt.Sprintf("action=input_issue_desc&serial=%s&category_id=%d", serialNumber, categoryID),
						"displayText": "พิมพ์รายละเอียดปัญหา",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "secondary", "margin": "sm",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "⏭️ ข้าม (ไม่ระบุ)",
						"data":        fmt.Sprintf("action=submit_issue&serial=%s&category_id=%d&desc=", serialNumber, categoryID),
						"displayText": "ข้ามการระบุรายละเอียด",
					},
				},
			},
		},
	}
}
