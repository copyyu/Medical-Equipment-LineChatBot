package templates

import (
	"fmt"

	"medical-webhook/internal/domain/line/entity"
)

// GetCategorySelectionFlex returns a Flex Message with category options as buttons
func GetCategorySelectionFlex(serialNumber string, categories []entity.TicketCategory) map[string]interface{} {
	// Build category buttons
	buttons := make([]interface{}, 0)
	for _, cat := range categories {
		buttons = append(buttons, map[string]interface{}{
			"type":   "button",
			"style":  "primary",
			"color":  cat.Color,
			"margin": "sm",
			"action": map[string]interface{}{
				"type":        "postback",
				"label":       fmt.Sprintf("%s %s", cat.Icon, cat.Name),
				"data":        fmt.Sprintf("action=confirm_category&serial=%s&category_id=%d", serialNumber, cat.ID),
				"displayText": fmt.Sprintf("เลือกหมวดหมู่: %s", cat.Name),
			},
		})
	}

	// If no categories, show default message
	if len(buttons) == 0 {
		buttons = append(buttons, map[string]interface{}{
			"type":  "button",
			"style": "primary",
			"color": "#78909C",
			"action": map[string]interface{}{
				"type":        "postback",
				"label":       "🔧 แจ้งซ่อมทั่วไป",
				"data":        fmt.Sprintf("action=confirm_category&serial=%s&category_id=0", serialNumber),
				"displayText": "เลือกหมวดหมู่: แจ้งซ่อมทั่วไป",
			},
		})
	}

	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#4A90D9",
			"paddingAll":      "12px",
			"contents": []interface{}{
				map[string]interface{}{"type": "text", "text": "📋 เลือกหมวดหมู่", "color": "#FFFFFF", "size": "md", "weight": "bold"},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("เลขเครื่อง: %s", serialNumber), "size": "sm", "color": "#888888",
				},
				map[string]interface{}{"type": "separator", "margin": "md"},
				map[string]interface{}{
					"type": "text", "text": "กรุณาเลือกประเภทการแจ้ง:", "size": "sm", "margin": "md", "wrap": true,
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": buttons,
		},
	}
}
