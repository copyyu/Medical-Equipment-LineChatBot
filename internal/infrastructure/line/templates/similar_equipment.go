package templates

import (
	"fmt"
	"medical-webhook/internal/domain/line/entity"
	"net/url"
)

// GetSimilarEquipmentFlex แสดงข้อมูลที่ใกล้เคียงที่สุด 1 รายการ พร้อม % ความใกล้เคียง
func GetSimilarEquipmentFlex(original string, bestIDCode string, similarityPct int) map[string]interface{} {
	escapedIDCode := url.QueryEscape(bestIDCode)

	// เลือกสี % ตามระดับความใกล้เคียง
	pctColor := "#4CAF50" // เขียว (สูง >= 70%)
	if similarityPct < 70 {
		pctColor = "#FF9800" // ส้ม (กลาง 50-69%)
	}
	if similarityPct < 50 {
		pctColor = "#F44336" // แดง (ต่ำ < 50%)
	}

	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#673AB7",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "🔎 ข้อมูลที่ใกล้เคียงที่สุด",
					"color":  "#FFFFFF",
					"weight": "bold",
				},
			},
		},
		"body": map[string]interface{}{
			"type":       "box",
			"layout":     "vertical",
			"spacing":    "md",
			"paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "text",
					"text":  "ระบบอ่านได้: " + original,
					"size":  "sm",
					"color": "#888888",
				},
				map[string]interface{}{
					"type":   "separator",
					"margin": "md",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   bestIDCode,
					"size":   "xxl",
					"weight": "bold",
					"color":  "#333333",
					"align":  "center",
					"margin": "md",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   fmt.Sprintf("ความใกล้เคียง %d%%", similarityPct),
					"size":   "md",
					"color":  pctColor,
					"weight": "bold",
					"align":  "center",
					"margin": "sm",
				},
			},
		},
		"footer": map[string]interface{}{
			"type":    "box",
			"layout":  "vertical",
			"spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "button",
					"style": "primary",
					"color": "#4CAF50",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "✅ ใช่ เลือกรายการนี้",
						"data":        "action=ocr_confirm_yes&serial=" + escapedIDCode,
						"displayText": "เลือก " + bestIDCode,
					},
				},
				map[string]interface{}{
					"type":  "button",
					"style": "link",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "❌ ไม่ใช่ (ถ่ายรูปใหม่)",
						"data":        "action=ocr_retake",
						"displayText": "ถ่ายรูปใหม่",
					},
				},
			},
		},
	}
}

// GetSimilarEquipmentListFlex แสดงรายการอุปกรณ์ที่ใกล้เคียง (แบบเดิม - หลายรายการ)
func GetSimilarEquipmentListFlex(original string, equipments []*entity.Equipment) map[string]interface{} {
	contents := []map[string]interface{}{}

	// Header
	contents = append(contents, map[string]interface{}{
		"type":   "text",
		"text":   "พบอุปกรณ์ที่ใกล้เคียง",
		"weight": "bold",
		"size":   "lg",
	})

	contents = append(contents, map[string]interface{}{
		"type":  "text",
		"text":  "จากข้อความ: " + original,
		"size":  "sm",
		"color": "#888888",
		"wrap":  true,
	})

	contents = append(contents, map[string]interface{}{
		"type":   "separator",
		"margin": "md",
	})

	limit := 5
	if len(equipments) < limit {
		limit = len(equipments)
	}

	for i := 0; i < limit; i++ {
		e := equipments[i]
		idCode := e.IDCode
		escapedIDCode := url.QueryEscape(idCode)

		contents = append(contents, map[string]interface{}{
			"type":   "button",
			"style":  "primary",
			"margin": "md",
			"action": map[string]interface{}{
				"type":  "postback",
				"label": idCode,
				"data":  "action=ocr_similar_select&serial=" + escapedIDCode + "&original=" + url.QueryEscape(original),
			},
		})
	}

	// ปุ่มถ่ายรูปใหม่
	contents = append(contents, map[string]interface{}{
		"type":   "button",
		"style":  "link",
		"margin": "md",
		"action": map[string]interface{}{
			"type":  "postback",
			"label": "ไม่มีข้อมูลที่ต้องการ (ถ่ายรูปใหม่)",
			"data":  "action=ocr_retake",
		},
	})

	return map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type":     "box",
			"layout":   "vertical",
			"contents": contents,
		},
	}
}

// GetSimilarConfirmFlex ถามยืนยันเมื่อผู้ใช้เลือกจากรายการใกล้เคียง
// แสดง: "ต้องการเปลี่ยนไปหมายเลข [selected] ใช่หรือไม่?" + "ที่ระบบอ่านได้คือ [original]"
func GetSimilarConfirmFlex(selectedIDCode string, originalOCR string) map[string]interface{} {
	escapedIDCode := url.QueryEscape(selectedIDCode)

	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#FF9800",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "⚠️ ยืนยันเปลี่ยนหมายเลข",
					"color":  "#FFFFFF",
					"weight": "bold",
				},
			},
		},
		"body": map[string]interface{}{
			"type":       "box",
			"layout":     "vertical",
			"spacing":    "md",
			"paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "text",
					"text":  "ที่ระบบอ่านได้คือ: " + originalOCR,
					"size":  "sm",
					"color": "#888888",
					"wrap":  true,
				},
				map[string]interface{}{
					"type":   "separator",
					"margin": "md",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "ต้องการเปลี่ยนไปหมายเลข",
					"size":   "md",
					"color":  "#333333",
					"align":  "center",
					"margin": "md",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   selectedIDCode,
					"size":   "xxl",
					"weight": "bold",
					"color":  "#1976D2",
					"align":  "center",
					"margin": "sm",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "ใช่หรือไม่?",
					"size":   "md",
					"color":  "#333333",
					"align":  "center",
					"margin": "sm",
				},
			},
		},
		"footer": map[string]interface{}{
			"type":    "box",
			"layout":  "horizontal",
			"spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "button",
					"style": "primary",
					"color": "#4CAF50",
					"flex":  1,
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "✅ ใช่",
						"data":        "action=ocr_confirm_yes&serial=" + escapedIDCode,
						"displayText": "เลือก " + selectedIDCode,
					},
				},
				map[string]interface{}{
					"type":  "button",
					"style": "primary",
					"color": "#F44336",
					"flex":  1,
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "❌ ไม่ใช่",
						"data":        "action=ocr_retake",
						"displayText": "ถ่ายรูปใหม่",
					},
				},
			},
		},
	}
}
