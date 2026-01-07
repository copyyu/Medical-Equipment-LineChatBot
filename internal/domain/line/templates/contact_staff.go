package templates

// GetContactStaffFlex returns a Flex Message for contact information
func GetContactStaffFlex() map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type": "box", "layout": "vertical", "backgroundColor": "#4CAF50", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "ติดต่อเจ้าหน้าที่", "color": "#FFFFFF", "size": "lg", "weight": "bold", "align": "center"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "ศูนย์เครื่องมือแพทย์", "weight": "bold", "size": "md"},
				map[string]interface{}{"type": "separator"},
				map[string]interface{}{"type": "box", "layout": "vertical", "margin": "md", "contents": []interface{}{
					map[string]interface{}{"type": "text", "text": "โทร: 123965845", "size": "sm"},
					map[string]interface{}{"type": "text", "text": "Email: lao@hospital.com", "size": "sm"},
					map[string]interface{}{"type": "text", "text": "เวลา: จ-ศ 08:00-17:00", "size": "sm"},
				}},
				map[string]interface{}{"type": "box", "layout": "vertical", "margin": "md", "backgroundColor": "#FFEBEE", "cornerRadius": "md", "paddingAll": "10px", "contents": []interface{}{
					map[string]interface{}{"type": "text", "text": "🚨 กรณีฉุกเฉิน", "weight": "bold", "color": "#C62828", "size": "sm"},
					map[string]interface{}{"type": "text", "text": "โทร: 12354675745 (24 ชม.)", "color": "#C62828", "size": "sm"},
				}},
			},
		},
	}
}
