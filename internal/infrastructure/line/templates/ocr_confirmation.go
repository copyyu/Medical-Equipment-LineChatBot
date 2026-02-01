package templates

import "fmt"

// GetOCRConfirmationFlex returns a Flex Message for OCR confirmation
// Shows detected serial number and asks user to confirm
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
				map[string]interface{}{"type": "text", "text": "🔍 ตรวจสอบหมายเลขเครื่อง", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "ระบบอ่านได้เลข:", "size": "sm", "color": "#888888",
				},
				map[string]interface{}{
					"type": "text", "text": detectedText, "size": "xxl", "weight": "bold", "color": "#0367D3", "align": "center",
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "text", "text": "ถูกต้องหรือไม่?", "size": "md", "weight": "bold", "align": "center", "margin": "md",
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "horizontal", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "primary", "color": "#4CAF50",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "✅ ใช่ ถูกต้อง",
						"data":        fmt.Sprintf("action=ocr_confirm_yes&serial=%s", detectedText),
						"displayText": "ใช่ ถูกต้อง",
					},
				},
				map[string]interface{}{
					"type": "button", "style": "secondary",
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       "❌ ไม่ใช่",
						"data":        "action=ocr_confirm_no",
						"displayText": "ไม่ใช่",
					},
				},
			},
		},
	}
}

// GetOCRNotFoundFlex returns a Flex Message when serial number is not found in DB
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
				map[string]interface{}{"type": "text", "text": "⚠️ ไม่พบข้อมูล", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("ระบบอ่านได้: %s", detectedText), "size": "md", "wrap": true,
				},
				map[string]interface{}{
					"type": "text", "text": "ไม่พบเครื่องมือนี้ในระบบ", "size": "sm", "color": "#888888", "wrap": true,
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "text", "text": "กรุณาตรวจสอบข้อความแล้วส่งรูปภาพใหม่", "size": "sm", "align": "center", "margin": "md",
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "กรุณาเลือกเมนูด้านล่างเพื่อทำรายการอื่น", "size": "xs", "color": "#888888", "align": "center",
				},
			},
		},
	}
}

// GetOCRErrorFlex returns a Flex Message when OCR fails
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
				map[string]interface{}{"type": "text", "text": "❌ อ่านรูปไม่สำเร็จ", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "ระบบไม่สามารถอ่านหมายเลขจากรูปได้", "size": "sm", "wrap": true,
				},
				map[string]interface{}{
					"type": "text", "text": "กรุณาถ่ายรูปใหม่ให้เห็นตัวเลขชัดๆ", "size": "sm", "color": "#888888", "wrap": true, "margin": "md",
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "กรุณาเลือกเมนูด้านล่างเพื่อทำรายการอื่น", "size": "xs", "color": "#888888", "align": "center",
				},
			},
		},
	}
}

// GetRetryPhotoFlex returns a message asking user to retake photo
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
				map[string]interface{}{"type": "text", "text": "📷 ส่งรูปใหม่", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "ขออภัยครับ รูปอาจไม่ชัด", "size": "md", "wrap": true,
				},
				map[string]interface{}{
					"type": "text", "text": "กรุณาถ่ายรูปใหม่ให้เห็นตัวเลขชัดๆ หน่อยครับ", "size": "sm", "color": "#888888", "wrap": true, "margin": "md",
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "หรือเลือกเมนูด้านล่างเพื่อทำรายการอื่น", "size": "xs", "color": "#888888", "align": "center",
				},
			},
		},
	}
}
