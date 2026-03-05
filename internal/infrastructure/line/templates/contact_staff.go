package templates

// GetContactStaffFlex returns a Flex Message for contact information
func GetContactStaffFlex() map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type": "box", "layout": "vertical", "backgroundColor": ColorPrimary, "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "📞 ติดต่อเจ้าหน้าที่", "color": ColorWhite, "size": "lg", "weight": "bold", "align": "center"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "🏥 ศูนย์เครื่องมือแพทย์", "weight": "bold", "size": "md", "color": ColorText},
				map[string]interface{}{"type": "separator"},
				map[string]interface{}{"type": "box", "layout": "vertical", "margin": "md", "contents": []interface{}{
					map[string]interface{}{"type": "text", "text": "📱 โทร: 0123456789", "size": "sm", "color": ColorText},
					map[string]interface{}{"type": "text", "text": "📧 Email: example@hospital.com", "size": "sm", "color": ColorText},
					map[string]interface{}{"type": "text", "text": "🕘 เวลา: จ-ศ 08:00-17:00", "size": "sm", "color": ColorTextLight},
				}},
				map[string]interface{}{"type": "box", "layout": "vertical", "margin": "md", "backgroundColor": ColorDangerLight, "cornerRadius": "md", "paddingAll": "10px", "contents": []interface{}{
					map[string]interface{}{"type": "text", "text": "🚨 กรณีฉุกเฉิน", "weight": "bold", "color": ColorDanger, "size": "sm"},
					map[string]interface{}{"type": "text", "text": "โทร: 9876543210 (24 ชม.)", "color": ColorDanger, "size": "sm"},
				}},
			},
		},
	}
}
