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
			"backgroundColor": ColorPrimary,
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "🔧 ข้อมูลเครื่องมือ", "color": ColorWhite, "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เลขเครื่อง: %s", serialNumber), "size": "md", "weight": "bold", "color": ColorAccent,
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "text", "text": "ต้องการดูข้อมูลอะไร?", "size": "sm", "margin": "md", "color": ColorText,
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "primary", "color": ColorAccent,
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "📋 ดูประวัติการซ่อม",
						"data":        fmt.Sprintf("action=view_repair_history&serial=%s", serialNumber),
						"displayText": "ดูประวัติการซ่อม",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "primary", "color": ColorWarning, "margin": "sm",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "📅 ดูอายุ/วงจรชีวิตเครื่อง",
						"data":        fmt.Sprintf("action=view_lifecycle&serial=%s", serialNumber),
						"displayText": "ดูอายุ/วงจรชีวิตเครื่อง",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "primary", "color": ColorInfo, "margin": "sm",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "📊 ดูสเปกเครื่อง",
						"data":        fmt.Sprintf("action=view_specs&serial=%s", serialNumber),
						"displayText": "ดูสเปกเครื่อง",
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
			"type": "text", "text": fmt.Sprintf("เครื่อง: %s", serialNumber), "size": "sm", "color": ColorTextLight,
		},
		map[string]interface{}{"type": "separator", "margin": "md"},
	}

	if len(records) == 0 {
		contents = append(contents, map[string]interface{}{
			"type": "text", "text": "ไม่พบประวัติการซ่อม", "size": "sm", "color": ColorTextLight, "margin": "md",
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
						"type": "text", "text": fmt.Sprintf("📅 %v", record["date"]), "size": "sm", "weight": "bold", "color": ColorText,
					},
					map[string]interface{}{
						"type": "text", "text": fmt.Sprintf("ประเภท: %v", record["type"]), "size": "xs", "color": ColorTextLight,
					},
					map[string]interface{}{
						"type": "text", "text": fmt.Sprintf("รายละเอียด: %v", record["description"]), "size": "xs", "color": ColorTextLight, "wrap": true,
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
			"backgroundColor": ColorAccent,
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "📋 ประวัติการซ่อม", "color": ColorWhite, "size": "md", "weight": "bold"},
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
			"backgroundColor": ColorWarning,
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "📅 อายุ/วงจรชีวิตเครื่อง", "color": ColorWhite, "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เครื่อง: %s", serialNumber), "size": "sm", "color": ColorTextLight,
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "md",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "อายุเครื่อง:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v ปี", data["equipment_age"]), "size": "sm", "flex": 1, "align": "end", "color": ColorText},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "อายุการใช้งานคาดหวัง:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v ปี", data["life_expectancy"]), "size": "sm", "flex": 1, "align": "end", "color": ColorText},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "อายุคงเหลือ:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v ปี", data["remain_life"]), "size": "sm", "flex": 1, "align": "end", "color": ColorText},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "% การใช้งาน:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%.1f%%", data["useful_percent"]), "size": "sm", "flex": 1, "align": "end", "color": ColorText},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "ปีที่ต้องเปลี่ยน:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["replacement_year"]), "size": "sm", "flex": 1, "align": "end", "color": ColorDanger},
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
			"backgroundColor": ColorInfo,
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "📊 สเปกเครื่อง", "color": ColorWhite, "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เครื่อง: %s", serialNumber), "size": "sm", "color": ColorTextLight,
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "md",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "รุ่น:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["model_name"]), "size": "sm", "flex": 2, "align": "end", "wrap": true, "color": ColorText},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "ยี่ห้อ:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["brand"]), "size": "sm", "flex": 2, "align": "end", "color": ColorText},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "แผนก:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["department"]), "size": "sm", "flex": 2, "align": "end", "wrap": true, "color": ColorText},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "วันที่รับ:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", data["receive_date"]), "size": "sm", "flex": 2, "align": "end", "color": ColorText},
					},
				},
				map[string]interface{}{
					"type": "box", "layout": "horizontal", "margin": "sm",
					"contents": []interface{}{
						map[string]interface{}{"type": "text", "text": "ราคา:", "size": "sm", "flex": 1, "color": ColorText},
						map[string]interface{}{"type": "text", "text": fmt.Sprintf("฿%.2f", data["price"]), "size": "sm", "flex": 2, "align": "end", "color": ColorText},
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
			},
		},
	}
}
