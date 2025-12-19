package templates

func GetMainMenuFlex() map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#0367D3",
			"paddingAll":      "15px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "ระบบเครื่องมือแพทย์", "color": "#FFFFFF", "size": "lg", "weight": "bold"},
				map[string]interface{}{"type": "text", "text": "Medical Equipment Service", "color": "#B8D4F0", "size": "xs"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "กรุณาเลือกบริการ", "weight": "bold", "size": "md"},
				map[string]interface{}{"type": "separator"},
				map[string]interface{}{"type": "button", "style": "primary", "color": "#5B9BD5", "margin": "md", "action": map[string]interface{}{"type": "postback", "label": "แจ้งปัญหา / เช็กสถานะ", "data": "action=report_problem", "displayText": "เมนูหลัก > แจ้งปัญหา / เช็กสถานะ"}},
				map[string]interface{}{"type": "button", "style": "primary", "color": "#FF9800", "margin": "sm", "action": map[string]interface{}{"type": "postback", "label": "🔄 แจ้งเปลี่ยนเครื่อง", "data": "action=request_change", "displayText": "เมนูหลัก > แจ้งเปลี่ยนเครื่อง"}},
				map[string]interface{}{"type": "button", "style": "primary", "color": "#4CAF50", "margin": "sm", "action": map[string]interface{}{"type": "postback", "label": "ติดตามสถานะ", "data": "action=track_status", "displayText": "เมนูหลัก > ติดตามสถานะ"}},
				map[string]interface{}{"type": "button", "style": "secondary", "margin": "sm", "action": map[string]interface{}{"type": "postback", "label": "ติดต่อเจ้าหน้าที่", "data": "action=contact_staff", "displayText": "เมนูหลัก > ติดต่อเจ้าหน้าที่"}},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "พิมพ์อะไรก็ได้เพื่อเรียกเมนูนี้", "size": "xxs", "color": "#AAAAAA", "align": "center"},
			},
		},
	}
}
