package handlers

// Flex Message templates for LINE Bot

// GetMainMenuFlex returns a Flex Message for the main menu
func GetMainMenuFlex() map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#0367D3",
			"paddingAll":      "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "ระบบเครื่องมือแพทย์",
					"color":  "#FFFFFF",
					"size":   "lg",
					"weight": "bold",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "Medical Equipment Service",
					"color": "#B8D4F0",
					"size":  "xs",
				},
			},
		},
		"body": map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"spacing": "md",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "กรุณาเลือกบริการ",
					"weight": "bold",
					"size":   "md",
				},
				map[string]interface{}{
					"type": "separator",
				},
				// แจ้งปัญหา / เช็กสถานะ
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"color":  "#5B9BD5",
					"margin": "md",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "แจ้งปัญหา / เช็กสถานะ",
						"data":  "action=report_problem",
					},
				},
				// แจ้งเปลี่ยนเครื่อง
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"color":  "#FF9800",
					"margin": "sm",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "🔄 แจ้งเปลี่ยนเครื่อง",
						"data":  "action=request_change",
					},
				},
				// ติดตามสถานะ
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"color":  "#4CAF50",
					"margin": "sm",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "ติดตามสถานะ",
						"data":  "action=track_status",
					},
				},
				// ติดต่อเจ้าหน้าที่
				map[string]interface{}{
					"type":   "button",
					"style":  "secondary",
					"margin": "sm",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "ติดต่อเจ้าหน้าที่",
						"data":  "action=contact_staff",
					},
				},
			},
		},
		"footer": map[string]interface{}{
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "text",
					"text":  "พิมพ์อะไรก็ได้เพื่อเรียกเมนูนี้",
					"size":  "xxs",
					"color": "#AAAAAA",
					"align": "center",
				},
			},
		},
	}
}

// GetEquipmentChangeFlex returns a Flex Message for equipment change request with a link
func GetEquipmentChangeFlex(linkURL string) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#FF9800",
			"paddingAll":      "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "🔄 แจ้งเปลี่ยนเครื่อง",
					"color":  "#FFFFFF",
					"size":   "lg",
					"weight": "bold",
					"align":  "center",
				},
			},
		},
		"body": map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"spacing": "md",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "กรุณากดปุ่มด้านล่างเพื่อเข้าสู่ระบบแจ้งเปลี่ยนเครื่อง",
					"wrap": true,
					"size": "sm",
				},
				map[string]interface{}{
					"type": "separator",
				},
				map[string]interface{}{
					"type":   "box",
					"layout": "vertical",
					"margin": "md",
					"contents": []interface{}{
						map[string]interface{}{
							"type":  "text",
							"text":  "-",
							"size":  "xs",
							"color": "#666666",
						},
						map[string]interface{}{
							"type":  "text",
							"text":  "-",
							"size":  "xs",
							"color": "#666666",
						},
						map[string]interface{}{
							"type":  "text",
							"text":  "✅ ยืนยันคำขอเปลี่ยนเครื่อง",
							"size":  "xs",
							"color": "#666666",
						},
					},
				},
			},
		},
		"footer": map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"color":  "#FF9800",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "เปิดฟอร์มแจ้งเปลี่ยน",
						"uri":   linkURL,
					},
				},
				map[string]interface{}{
					"type":   "button",
					"style":  "secondary",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "↩️ กลับเมนูหลัก",
						"data":  "action=main_menu",
					},
				},
			},
		},
	}
}

// GetContactStaffFlex returns a Flex Message for contact information
func GetContactStaffFlex() map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#4CAF50",
			"paddingAll":      "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "ติดต่อเจ้าหน้าที่",
					"color":  "#FFFFFF",
					"size":   "lg",
					"weight": "bold",
					"align":  "center",
				},
			},
		},
		"body": map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"spacing": "md",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "ศูนย์เครื่องมือแพทย์",
					"weight": "bold",
					"size":   "md",
				},
				map[string]interface{}{
					"type": "separator",
				},
				map[string]interface{}{
					"type":   "box",
					"layout": "vertical",
					"margin": "md",
					"contents": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": "โทร: 123965845",
							"size": "sm",
						},
						map[string]interface{}{
							"type": "text",
							"text": "Email: lao@hospital.com",
							"size": "sm",
						},
						map[string]interface{}{
							"type": "text",
							"text": "เวลา: จ-ศ 08:00-17:00",
							"size": "sm",
						},
					},
				},
				map[string]interface{}{
					"type":            "box",
					"layout":          "vertical",
					"margin":          "md",
					"backgroundColor": "#FFEBEE",
					"cornerRadius":    "md",
					"paddingAll":      "10px",
					"contents": []interface{}{
						map[string]interface{}{
							"type":   "text",
							"text":   "🚨 กรณีฉุกเฉิน",
							"weight": "bold",
							"color":  "#C62828",
							"size":   "sm",
						},
						map[string]interface{}{
							"type":  "text",
							"text":  "โทร: 12354675745 (24 ชม.)",
							"color": "#C62828",
							"size":  "sm",
						},
					},
				},
			},
		},
		"footer": map[string]interface{}{
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "secondary",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "postback",
						"label": "↩️ กลับเมนูหลัก",
						"data":  "action=main_menu",
					},
				},
			},
		},
	}
}
