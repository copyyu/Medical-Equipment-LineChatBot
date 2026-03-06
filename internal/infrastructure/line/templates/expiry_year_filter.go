package templates

import (
	"fmt"
	"time"
)

// GetExpiryYearFilterFlex returns a Flex Message asking user to choose year filter
// Shows summary counts and 3 buttons: this year / next year / all
func GetExpiryYearFilterFlex(deptName string, departmentID uint, expiredCount int, nearExpiryCount int) map[string]interface{} {
	thisYear := time.Now().Year()
	nextYear := thisYear + 1
	totalCount := expiredCount + nearExpiryCount

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
					"type": "text", "text": "📊 เครื่องใกล้หมดอายุ",
					"color": ColorWhite, "size": "lg", "weight": "bold",
				},
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("แผนก: %s", deptName),
					"color": "#FFFFFFCC", "size": "sm",
				},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "md", "paddingAll": "15px",
			"contents": []interface{}{
				// Summary count section
				map[string]interface{}{
					"type": "box", "layout": "vertical", "spacing": "sm",
					"backgroundColor": ColorBgAlt, "cornerRadius": "8px", "paddingAll": "12px",
					"contents": []interface{}{
						map[string]interface{}{
							"type": "box", "layout": "horizontal",
							"contents": []interface{}{
								map[string]interface{}{"type": "text", "text": fmt.Sprintf("🔴 หมดอายุปี %d:", thisYear), "size": "sm", "color": ColorDanger, "flex": 4},
								map[string]interface{}{"type": "text", "text": fmt.Sprintf("%d เครื่อง", expiredCount), "size": "sm", "color": ColorDanger, "weight": "bold", "align": "end", "flex": 2},
							},
						},
						map[string]interface{}{
							"type": "box", "layout": "horizontal", "margin": "sm",
							"contents": []interface{}{
								map[string]interface{}{"type": "text", "text": fmt.Sprintf("🟡 หมดอายุปี %d:", nextYear), "size": "sm", "color": ColorWarning, "flex": 4},
								map[string]interface{}{"type": "text", "text": fmt.Sprintf("%d เครื่อง", nearExpiryCount), "size": "sm", "color": ColorWarning, "weight": "bold", "align": "end", "flex": 2},
							},
						},
						map[string]interface{}{"type": "separator", "margin": "md"},
						map[string]interface{}{
							"type": "box", "layout": "horizontal", "margin": "md",
							"contents": []interface{}{
								map[string]interface{}{"type": "text", "text": "📋 รวมทั้งหมด:", "size": "sm", "color": ColorText, "weight": "bold", "flex": 4},
								map[string]interface{}{"type": "text", "text": fmt.Sprintf("%d เครื่อง", totalCount), "size": "sm", "color": ColorText, "weight": "bold", "align": "end", "flex": 2},
							},
						},
					},
				},
				map[string]interface{}{
					"type": "text", "text": "ต้องการดูข้อมูลช่วงไหน?",
					"size": "md", "weight": "bold", "margin": "lg", "align": "center", "color": ColorText,
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box", "layout": "vertical", "spacing": "sm",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "button", "style": "primary", "color": ColorDanger,
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       fmt.Sprintf("🔴 ดูเฉพาะปี %d (%d)", thisYear, expiredCount),
						"data":        fmt.Sprintf("action=view_expiry_filtered&department_id=%d&filter=this_year", departmentID),
						"displayText": fmt.Sprintf("ดูเครื่องหมดอายุปี %d", thisYear),
					},
				},
				map[string]interface{}{
					"type": "button", "style": "primary", "color": ColorWarning,
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       fmt.Sprintf("🟡 ดูเฉพาะปี %d (%d)", nextYear, nearExpiryCount),
						"data":        fmt.Sprintf("action=view_expiry_filtered&department_id=%d&filter=next_year", departmentID),
						"displayText": fmt.Sprintf("ดูเครื่องหมดอายุปี %d", nextYear),
					},
				},
				map[string]interface{}{
					"type": "button", "style": "primary", "color": ColorPrimaryDark,
					"action": map[string]interface{}{
						"type":        "postback",
						"label":       fmt.Sprintf("📋 ดูทั้งหมด (%d)", totalCount),
						"data":        fmt.Sprintf("action=view_expiry_filtered&department_id=%d&filter=all", departmentID),
						"displayText": "ดูทั้งหมด",
					},
				},
			},
		},
	}
}
