package templates

import "fmt"

//
// ===============================
// OCR CONFIRMATION
// ===============================

// GetOCRConfirmationFlex
// แสดงเลขที่ OCR อ่านได้ และให้ผู้ใช้ยืนยัน
func GetOCRConfirmationFlex(detectedText string, imageURL string) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#0367D3",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "🔍 ตรวจสอบหมายเลขเครื่อง",
					"color":  "#FFFFFF",
					"size":   "md",
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
					"text":  "ระบบอ่านได้เลข:",
					"size":  "sm",
					"color": "#888888",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   detectedText,
					"size":   "xxl",
					"weight": "bold",
					"color":  "#0367D3",
					"align":  "center",
				},
				map[string]interface{}{
					"type":   "separator",
					"margin": "md",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "ถูกต้องหรือไม่?",
					"size":   "md",
					"weight": "bold",
					"align":  "center",
					"margin": "md",
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
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "✅ ใช่ ถูกต้อง",
						"data":        fmt.Sprintf("action=ocr_confirm_yes&serial=%s", detectedText),
						"displayText": "ใช่ ถูกต้อง",
					},
				},
				map[string]interface{}{
					"type":  "button",
					"style": "secondary",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "❌ ไม่ใช่",
						"data":        fmt.Sprintf("action=ocr_confirm_no&serial=%s", detectedText),
						"displayText": "ไม่ใช่",
					},
				},
			},
		},
	}
}

//
// ===============================
// OCR SIMILAR SERIALS (NEW)
// ===============================

// GetOCRSimilarFlex
// แสดงหมายเลขที่ใกล้เคียงจาก database เมื่อผู้ใช้กด "ไม่ใช่"
func GetOCRSimilarFlex(detectedText string, suggestions []string) map[string]interface{} {

	contents := []interface{}{
		map[string]interface{}{
			"type":  "text",
			"text":  fmt.Sprintf("ระบบอ่านได้: %s", detectedText),
			"size":  "sm",
			"color": "#888888",
		},
		map[string]interface{}{
			"type":   "text",
			"text":   "หมายเลขที่ต้องการใช่เหล่านี้หรือไม่",
			"size":   "md",
			"weight": "bold",
			"margin": "md",
		},
	}

	for _, s := range suggestions {
		contents = append(contents, map[string]interface{}{
			"type":   "button",
			"style":  "secondary",
			"margin": "sm",
			"action": map[string]interface{}{
				"type":        "postback",
				"label":       s,
				"data":        fmt.Sprintf("action=ocr_select_serial&serial=%s", s),
				"displayText": s,
			},
		})
	}

	contents = append(contents, map[string]interface{}{
		"type":   "button",
		"style":  "secondary",
		"margin": "md",
		"action": map[string]interface{}{
			"type":        "postback",
			"label":       "❓ ไม่มีในรายการ",
			"data":        "action=ocr_no_match",
			"displayText": "ไม่มีในรายการ",
		},
	})

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
					"text":   "🔎 หมายเลขที่ใกล้เคียง",
					"color":  "#FFFFFF",
					"weight": "bold",
				},
			},
		},
		"body": map[string]interface{}{
			"type":     "box",
			"layout":   "vertical",
			"spacing":  "md",
			"contents": contents,
		},
	}
}

//
// ===============================
// OCR NOT FOUND
// ===============================

// GetOCRNotFoundFlex
func GetOCRNotFoundFlex(detectedText string) map[string]interface{} {
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
					"text":   "⚠️ ไม่พบในฐานระบบ",
					"color":  "#FFFFFF",
					"weight": "bold",
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
					"text":   detectedText,
					"size":   "xl",
					"weight": "bold",
					"align":  "center",
					"color":  "#333333",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "ไม่พบในฐานระบบ หรือภาพไม่ชัด กรุณาส่งภาพมาใหม่",
					"size":  "sm",
					"wrap":  true,
					"color": "#888888",
				},
			},
		},
	}
}

//
// ===============================
// OCR ERROR
// ===============================

// GetOCRErrorFlex
func GetOCRErrorFlex() map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#F44336",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "❌ อ่านรูปไม่สำเร็จ",
					"color":  "#FFFFFF",
					"weight": "bold",
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
					"text": "ระบบไม่สามารถอ่านหมายเลขจากรูปได้",
					"size": "sm",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "กรุณาถ่ายรูปใหม่ให้เห็นตัวเลขชัดๆ",
					"size":   "sm",
					"color":  "#888888",
					"margin": "md",
				},
			},
		},
	}
}

//
// ===============================
// RETRY PHOTO
// ===============================

// GetRetryPhotoFlex
func GetRetryPhotoFlex() map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#5B9BD5",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "📷 ส่งรูปใหม่",
					"color":  "#FFFFFF",
					"weight": "bold",
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
					"text": "ขออภัยครับ รูปอาจไม่ชัด",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "กรุณาถ่ายรูปใหม่ให้เห็นตัวเลขชัดๆ หน่อยครับ",
					"size":   "sm",
					"color":  "#888888",
					"margin": "md",
				},
			},
		},
	}
}
