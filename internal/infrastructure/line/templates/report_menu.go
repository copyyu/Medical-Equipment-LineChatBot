package templates

// GetReportMenuFlex returns a Flex Message sub-menu for "แจ้งปัญหา / เช็คสถานะ"
// with options: report problem or view equipment expiry
func GetReportMenuFlex() map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#1B5E20",
			"paddingAll":      "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "🔧 แจ้งปัญหา / เช็คสถานะ",
					"color": "#FFFFFF", "size": "lg", "weight": "bold",
				},
				map[string]interface{}{
					"type": "text", "text": "เลือกบริการที่ต้องการค่ะ",
					"color": "#FFFFFFCC", "size": "sm",
				},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "primary", "color": "#FF5722", "height": "md",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "📸 แจ้งปัญหา / เช็คสถานะเครื่อง",
						"data":        "action=start_report_mode",
						"displayText": "แจ้งปัญหา / เช็คสถานะเครื่อง",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "primary", "color": "#FF9800", "height": "md",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "📊 ดูเครื่องใกล้หมดอายุ",
						"data":        "action=view_equipment_expiry",
						"displayText": "ดูเครื่องใกล้หมดอายุ",
					},
				},
			},
		},
	}
}
