package templates

import (
	"fmt"

	"medical-webhook/internal/domain/line/entity"
)

// GetDepartmentSelectionFlex returns a Flex Message for selecting department before viewing expiring equipment
func GetDepartmentSelectionFlex(departments []entity.Department) map[string]interface{} {
	// Build department buttons
	deptButtons := []interface{}{}

	for _, dept := range departments {
		deptButtons = append(deptButtons, map[string]interface{}{
			"type":   "button",
			"style":  "secondary",
			"color":  "#4CAF50",
			"height": "sm",
			"action": map[string]interface{}{
				"type":        "postback",
				"label":       dept.Name,
				"data":        fmt.Sprintf("action=view_equipment_expiry_by_dept&department_id=%d", dept.ID),
				"displayText": "เลือก " + dept.Name,
			},
		})
	}

	// Limit to 10 departments for display
	if len(deptButtons) > 10 {
		deptButtons = deptButtons[:10]
	}

	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#1B5E20",
			"paddingAll":      "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "🏥 เลือกแผนก",
					"color": "#FFFFFF", "size": "lg", "weight": "bold",
				},
				map[string]interface{}{
					"type": "text", "text": "เลือกแผนกเพื่อดูเครื่องใกล้หมดอายุ",
					"color": "#FFFFFFCC", "size": "sm",
				},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm", "paddingAll": "15px",
			"contents": deptButtons,
		},
	}
}

// GetDepartmentSelectionWithInputFlex returns Flex Message with department list + manual input option
func GetDepartmentSelectionWithInputFlex(departments []entity.Department) map[string]interface{} {
	// Build department buttons (limit to 6 for space)
	deptButtons := []interface{}{}

	for i, dept := range departments {
		if i >= 6 {
			break
		}
		deptButtons = append(deptButtons, map[string]interface{}{
			"type":   "button",
			"style":  "secondary",
			"color":  "#4CAF50",
			"height": "sm",
			"action": map[string]interface{}{
				"type":        "postback",
				"label":       dept.Name,
				"data":        fmt.Sprintf("action=view_equipment_expiry_by_dept&department_id=%d", dept.ID),
				"displayText": "เลือก " + dept.Name,
			},
		})
	}

	// Add separator and input instruction
	contents := []interface{}{}
	contents = append(contents, map[string]interface{}{
		"type":  "text",
		"text":  "เลือกแผนกจากรายการด้านล่าง หรือพิมพ์ชื่อแผนกของคุณ",
		"size":  "xs",
		"color": "#666666",
		"wrap":  true,
	})
	contents = append(contents, map[string]interface{}{
		"type":   "separator",
		"margin": "md",
	})
	contents = append(contents, deptButtons...)

	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": "#1B5E20",
			"paddingAll":      "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "🏥 เลือกแผนกของคุณ",
					"color": "#FFFFFF", "size": "lg", "weight": "bold",
				},
				map[string]interface{}{
					"type": "text", "text": "เพื่อดูเครื่องมือใกล้หมดอายุในแผนก",
					"color": "#FFFFFFCC", "size": "sm",
				},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm", "paddingAll": "15px",
			"contents": contents,
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm", "paddingAll": "10px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "text",
					"text":  "💡 หากไม่พบแผนกของคุณ พิมพ์ชื่อแผนกได้เลยค่ะ",
					"size":  "xxs",
					"color": "#888888",
					"align": "center",
				},
			},
		},
	}
}
