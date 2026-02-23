package templates

import (
	"fmt"

	"medical-webhook/internal/domain/line/entity"
)

// deptButtonsPerBubble defines how many department buttons fit in one bubble body
const deptButtonsPerBubble = 6

// GetDepartmentSelectionFlex returns a single-bubble Flex Message for selecting department
func GetDepartmentSelectionFlex(departments []entity.Department) map[string]interface{} {
	deptButtons := buildDeptButtons(departments)

	// Limit to 10 departments for single bubble display
	if len(deptButtons) > 10 {
		deptButtons = deptButtons[:10]
	}

	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": ColorPrimaryDark,
			"paddingAll":      "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "🏥 เลือกแผนก",
					"color": ColorWhite, "size": "lg", "weight": "bold",
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

// GetDepartmentSelectionWithInputFlex returns Carousel Flex Message with department buttons
// split across multiple bubbles. If departments <= 6, returns a single bubble.
// Supports both button selection and text input.
func GetDepartmentSelectionWithInputFlex(departments []entity.Department) map[string]interface{} {
	allButtons := buildDeptButtons(departments)

	// ถ้าแผนกน้อย ใช้ bubble เดียว
	if len(allButtons) <= deptButtonsPerBubble {
		return buildSingleBubbleWithInput(allButtons)
	}

	// แผนกเยอะ → ใช้ Carousel แบ่ง bubble ละ deptButtonsPerBubble ปุ่ม
	return buildCarouselWithInput(allButtons, departments)
}

// deptColorPalette — สีแยกแผนก สุ่มวนไม่ซ้ำ (10 สี)
var deptColorPalette = []string{
	"#00897B", // Teal
	"#5C6BC0", // Indigo
	"#43A047", // Green
	"#F4511E", // Deep Orange
	"#1E88E5", // Blue
	"#8E24AA", // Purple
	"#00ACC1", // Cyan
	"#FB8C00", // Orange
	"#3949AB", // Dark Indigo
	"#D81B60", // Pink
}

// buildDeptButtons creates postback button elements for each department
func buildDeptButtons(departments []entity.Department) []interface{} {
	buttons := []interface{}{}
	for i, dept := range departments {
		color := deptColorPalette[i%len(deptColorPalette)]
		buttons = append(buttons, map[string]interface{}{
			"type":   "button",
			"style":  "primary",
			"color":  color,
			"height": "sm",
			"action": map[string]interface{}{
				"type":        "postback",
				"label":       dept.Name,
				"data":        fmt.Sprintf("action=view_equipment_expiry_by_dept&department_id=%d", dept.ID),
				"displayText": "เลือก " + dept.Name,
			},
		})
	}
	return buttons
}

// buildSingleBubbleWithInput creates a single bubble with buttons + text input instruction
func buildSingleBubbleWithInput(buttons []interface{}) map[string]interface{} {
	contents := []interface{}{
		map[string]interface{}{
			"type":  "text",
			"text":  "เลือกแผนกจากรายการด้านล่าง หรือพิมพ์ชื่อแผนก",
			"size":  "xs",
			"color": ColorTextLight,
			"wrap":  true,
		},
		map[string]interface{}{
			"type":   "separator",
			"margin": "md",
		},
	}
	contents = append(contents, buttons...)

	return map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
		"header": map[string]interface{}{
			"type":            "box",
			"layout":          "vertical",
			"backgroundColor": ColorPrimaryDark,
			"paddingAll":      "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "🏥 เลือกแผนกของคุณ",
					"color": ColorWhite, "size": "lg", "weight": "bold",
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
					"color": ColorTextLight,
					"align": "center",
					"wrap":  true,
				},
			},
		},
	}
}

// buildCarouselWithInput creates a carousel of bubbles, each containing up to deptButtonsPerBubble buttons
func buildCarouselWithInput(allButtons []interface{}, departments []entity.Department) map[string]interface{} {
	bubbles := []interface{}{}
	totalButtons := len(allButtons)

	for i := 0; i < totalButtons; i += deptButtonsPerBubble {
		end := i + deptButtonsPerBubble
		if end > totalButtons {
			end = totalButtons
		}
		chunk := allButtons[i:end]
		isFirst := (i == 0)
		isLast := (end >= totalButtons)

		bubble := buildCarouselBubble(chunk, isFirst, isLast, i/deptButtonsPerBubble+1, (totalButtons+deptButtonsPerBubble-1)/deptButtonsPerBubble)
		bubbles = append(bubbles, bubble)
	}

	// LINE Carousel supports up to 12 bubbles
	if len(bubbles) > 12 {
		bubbles = bubbles[:12]
	}

	return map[string]interface{}{
		"type":     "carousel",
		"contents": bubbles,
	}
}

// buildCarouselBubble creates one bubble in the carousel
func buildCarouselBubble(buttons []interface{}, isFirst, isLast bool, page, totalPages int) map[string]interface{} {
	bubble := map[string]interface{}{
		"type": "bubble",
		"size": "kilo",
	}

	// Header
	headerText := fmt.Sprintf("🏥 เลือกแผนก (%d/%d)", page, totalPages)
	headerContents := []interface{}{
		map[string]interface{}{
			"type": "text", "text": headerText,
			"color": ColorWhite, "size": "lg", "weight": "bold",
		},
	}

	if isFirst {
		headerContents = append(headerContents, map[string]interface{}{
			"type": "text", "text": "เลื่อน ← → เพื่อดูแผนกเพิ่มเติม",
			"color": "#FFFFFFCC", "size": "xs",
		})
	}

	bubble["header"] = map[string]interface{}{
		"type":            "box",
		"layout":          "vertical",
		"backgroundColor": ColorPrimaryDark,
		"paddingAll":      "15px",
		"contents":        headerContents,
	}

	// Body
	bodyContents := []interface{}{}
	if isFirst {
		bodyContents = append(bodyContents,
			map[string]interface{}{
				"type":  "text",
				"text":  "กดเลือกแผนก หรือพิมพ์ชื่อแผนกได้เลยค่ะ",
				"size":  "xs",
				"color": ColorTextLight,
				"wrap":  true,
			},
			map[string]interface{}{
				"type":   "separator",
				"margin": "md",
			},
		)
	}
	bodyContents = append(bodyContents, buttons...)

	bubble["body"] = map[string]interface{}{
		"type": "box", "layout": "vertical", "spacing": "sm", "paddingAll": "15px",
		"contents": bodyContents,
	}

	// Footer on last bubble
	if isLast {
		bubble["footer"] = map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm", "paddingAll": "10px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "text",
					"text":  "💡 หากไม่พบแผนกของคุณ พิมพ์ชื่อแผนกได้เลยค่ะ",
					"size":  "xxs",
					"color": ColorTextLight,
					"align": "center",
					"wrap":  true,
				},
			},
		}
	}

	return bubble
}
