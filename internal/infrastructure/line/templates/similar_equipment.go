package templates

import (
	"medical-webhook/internal/domain/line/entity"
	"net/url"
)

func GetSimilarEquipmentFlex(original string, equipments []*entity.Equipment) map[string]interface{} {
	contents := []map[string]interface{}{}

	// 🔹 Header
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

	// 🔹 จำกัดจำนวน (กัน Flex พัง)
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
				"data":  "action=ocr_confirm_yes&serial=" + escapedIDCode,
			},
		})
	}

	// 🔹 ปุ่มไม่มีข้อมูลที่ต้องการ (ถ่ายรูปใหม่)
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
