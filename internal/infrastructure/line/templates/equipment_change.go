package templates

// GetEquipmentChangeFlex returns a Flex Message for equipment change request
func GetEquipmentChangeFlex(linkURL string) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type": "box", "layout": "vertical", "backgroundColor": ColorWarning, "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "🔄 แจ้งเปลี่ยนเครื่อง", "color": ColorWhite, "size": "lg", "weight": "bold", "align": "center"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "กรุณากดปุ่มด้านล่างเพื่อเข้าสู่ระบบแจ้งเปลี่ยนเครื่อง", "wrap": true, "size": "sm", "color": ColorText},
				map[string]interface{}{"type": "separator"},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{"type": "button", "style": "primary", "color": ColorWarning, "height": "sm", "action": map[string]interface{}{"type": "uri", "label": "เปิดฟอร์มแจ้งเปลี่ยน", "uri": linkURL}},
			},
		},
	}
}
