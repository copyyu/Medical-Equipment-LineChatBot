package templates

import "fmt"

// GetEquipmentOptionsFlex returns a Flex Message with equipment action options
func GetEquipmentOptionsFlex(serialNumber string) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#4CAF50",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "🔧 ข้อมูลเครื่องมือ", "color": "#FFFFFF", "size": "md", "weight": "bold"},
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
					"type": "text", "text": "ต้องการดูข้อมูลอะไร?", "size": "sm", "margin": "md",
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
						"label":       "📋 ดูประวัติการซ่อม",
						"data":        fmt.Sprintf("action=view_repair_history&serial=%s", serialNumber),
						"displayText": "ดูประวัติการซ่อม",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "primary", "color": "#FF9800", "margin": "sm",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "📅 ดูอายุ/วงจรชีวิตเครื่อง",
						"data":        fmt.Sprintf("action=view_lifecycle&serial=%s", serialNumber),
						"displayText": "ดูอายุ/วงจรชีวิตเครื่อง",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "primary", "color": "#9C27B0", "margin": "sm",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "📊 ดูสเปกเครื่อง",
						"data":        fmt.Sprintf("action=view_specs&serial=%s", serialNumber),
						"displayText": "ดูสเปกเครื่อง",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "secondary", "margin": "md",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "🏠 กลับเมนูหลัก",
						"data":        "action=main_menu",
						"displayText": "กลับเมนูหลัก",
					},
				},
			},
		},
	}
}

// GetRepairHistoryFlex returns repair history as Flex Message
func GetRepairHistoryFlex(serialNumber string, records []map[string]interface{}) map[string]interface{} {
	contents := []interface{}{
		map[string]interface{}{
			"type": "text", "text": fmt.Sprintf("เครื่อง: %s", serialNumber), "size": "sm", "color": "#888888",
		},
		map[string]interface{}{"type": "separator", "margin": "md"},
	}

	if len(records) == 0 {
		contents = append(contents, map[string]interface{}{
			"type": "text", "text": "ไม่พบประวัติการซ่อม", "size": "sm", "color": "#888888", "margin": "md",
		})
	} else {
		for i, record := range records {
			if i >= 5 { // Limit to 5 records
				break
			}
			contents = append(contents, map[string]interface{}{
				"type": "box", "layout": "vertical", "margin": "md",
				"contents": []interface{}{
					map[string]interface{}{
						"type": "text", "text": fmt.Sprintf("📅 %v", record["date"]), "size": "sm", "weight": "bold",
					},
					map[string]interface{}{
						"type": "text", "text": fmt.Sprintf("ประเภท: %v", record["type"]), "size": "xs", "color": "#888888",
					},
					map[string]interface{}{
						"type": "text", "text": fmt.Sprintf("รายละเอียด: %v", record["description"]), "size": "xs", "color": "#888888", "wrap": true,
					},
				},
			})
		}
	}

	return map[string]interface{}{
		"type": "bubble",
		"size": "mega",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#5B9BD5",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "📋 ประวัติการซ่อม", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm", "paddingAll": "15px",
			"contents": contents,
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "horizontal", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "secondary",
					"action": map[string]interface{}{
						"type": "postback", "label": "⬅️ ย้อนกลับ",
						"data":        fmt.Sprintf("action=ocr_confirm_yes&serial=%s", serialNumber),
						"displayText": "ย้อนกลับ",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "secondary",
					"action": map[string]interface{}{
						"type": "postback", "label": "🏠 เมนูหลัก",
						"data": "action=main_menu", "displayText": "กลับเมนูหลัก",
					},
				},
			},
		},
	}
}

// GetLifecycleFlex returns equipment lifecycle info as Flex Message
func GetLifecycleFlex(serialNumber string, data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "mega",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#FF9800",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "📅 อายุ/วงจรชีวิตเครื่อง", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เครื่อง: %s", serialNumber), "size": "sm", "color": "#888888",
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "md",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "อายุเครื่อง:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v ปี", data["equipment_age"]), "size": "sm", "flex": 1, "align": "end"},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "อายุการใช้งานคาดหวัง:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v ปี", data["life_expectancy"]), "size": "sm", "flex": 1, "align": "end"},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "อายุคงเหลือ:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v ปี", data["remain_life"]), "size": "sm", "flex": 1, "align": "end"},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "% การใช้งาน:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%.1f%%", data["useful_percent"]), "size": "sm", "flex": 1, "align": "end"},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "ปีที่ต้องเปลี่ยน:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["replacement_year"]), "size": "sm", "flex": 1, "align": "end", "color": "#F44336"},
					},
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "horizontal", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "secondary",
					"action": map[string]interface{}{
						"type": "postback", "label": "⬅️ ย้อนกลับ",
						"data":        fmt.Sprintf("action=ocr_confirm_yes&serial=%s", serialNumber),
						"displayText": "ย้อนกลับ",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "secondary",
					"action": map[string]interface{}{
						"type": "postback", "label": "🏠 เมนูหลัก",
						"data": "action=main_menu", "displayText": "กลับเมนูหลัก",
					},
				},
			},
		},
	}
}

// GetSpecsFlex returns equipment specs as Flex Message
func GetSpecsFlex(serialNumber string, data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "mega",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#9C27B0",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "📊 สเปกเครื่อง", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เครื่อง: %s", serialNumber), "size": "sm", "color": "#888888",
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "md",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "รุ่น:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["model_name"]), "size": "sm", "flex": 2, "align": "end", "wrap": true},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "ยี่ห้อ:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["brand"]), "size": "sm", "flex": 2, "align": "end"},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "แผนก:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["department"]), "size": "sm", "flex": 2, "align": "end", "wrap": true},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "วันที่รับ:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["receive_date"]), "size": "sm", "flex": 2, "align": "end"},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "ราคา:", "size": "sm", "flex": 1},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("฿%.2f", data["price"]), "size": "sm", "flex": 2, "align": "end"},
					},
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "horizontal", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "secondary",
					"action": map[string]interface{}{
						"type": "postback", "label": "⬅️ ย้อนกลับ",
						"data":        fmt.Sprintf("action=ocr_confirm_yes&serial=%s", serialNumber),
						"displayText": "ย้อนกลับ",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "secondary",
					"action": map[string]interface{}{
						"type": "postback", "label": "🏠 เมนูหลัก",
						"data": "action=main_menu", "displayText": "กลับเมนูหลัก",
					},
				},
			},
		},
	}
}
