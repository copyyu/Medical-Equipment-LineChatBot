package templates

import (
	"fmt"
	"medical-webhook/internal/domain/line/entity"
	"time"
)

// GetEquipmentExpiryFlex returns a single Flex Message bubble with 2 sections stacked vertically:
func GetEquipmentExpiryFlex(expired []entity.Equipment, nearExpiry []entity.Equipment) map[string]interface{} {
	thisYear := time.Now().Year()
	nextYear := thisYear + 1

	bodyContents := []interface{}{}

	// Section 1: หมดอายุภายในปีนี้
	bodyContents = append(bodyContents, buildSectionHeader(
		fmt.Sprintf("🔴 หมดอายุภายในปี %d", thisYear), "#E53935"))
	bodyContents = append(bodyContents, buildEquipmentList(expired)...)
	bodyContents = append(bodyContents, map[string]interface{}{"type": "separator", "margin": "lg"})

	// Section 2: หมดอายุปีหน้า
	bodyContents = append(bodyContents, buildSectionHeader(
		fmt.Sprintf("🟡 หมดอายุปี %d", nextYear), "#FF9800"))
	bodyContents = append(bodyContents, buildEquipmentList(nearExpiry)...)

	return map[string]interface{}{
		"type": "bubble", "size": "mega",
		"header": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"backgroundColor": "#1B5E20", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "📋 รายการเครื่องใกล้หมดอายุ",
					"color": "#FFFFFF", "size": "lg", "weight": "bold",
				},
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("ข้อมูล ณ %s", time.Now().Format("02/01/2006")),
					"color": "#FFFFFFCC", "size": "xs",
				},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"spacing": "sm", "paddingAll": "12px",
			"contents": bodyContents,
		},
	}
}

func calcMonthsRemaining(replacementYear int) int {
	now := time.Now()
	return (replacementYear-now.Year())*12 - int(now.Month()) + 8
}

// buildSectionHeader สร้าง header ของแต่ละ section
func buildSectionHeader(title string, color string) map[string]interface{} {
	return map[string]interface{}{
		"type": "box", "layout": "horizontal", "margin": "lg",
		"paddingAll": "8px", "backgroundColor": color, "cornerRadius": "6px",
		"contents": []interface{}{
			map[string]interface{}{
				"type": "text", "text": title,
				"color": "#FFFFFF", "size": "sm", "weight": "bold",
			},
		},
	}
}

// buildEquipmentList สร้างรายการเครื่องมือ (running number + id + model + เดือน)
func buildEquipmentList(equipments []entity.Equipment) []interface{} {
	items := []interface{}{}

	if len(equipments) == 0 {
		items = append(items, map[string]interface{}{
			"type": "text", "text": "✅ ไม่มีเครื่องในหมวดนี้",
			"size": "xs", "color": "#4CAF50", "align": "center", "margin": "sm",
		})
		return items
	}

	items = append(items, map[string]interface{}{
		"type": "text", "text": fmt.Sprintf("📊 รวม %d รายการ", len(equipments)),
		"size": "xs", "color": "#888888", "margin": "sm",
	})

	for i, e := range equipments {
		items = append(items, buildEquipmentRow(i+1, e))
	}

	return items
}

// buildEquipmentRow — แต่ละรายการแสดง: running number + ID code + ชื่อเครื่อง + เหลือกี่เดือน
func buildEquipmentRow(num int, e entity.Equipment) map[string]interface{} {
	modelName := "-"
	if e.Model.ModelName != "" {
		modelName = e.Model.ModelName
	}

	monthsText := "-"
	monthsColor := "#888888"
	if e.ReplacementYear != nil {
		m := calcMonthsRemaining(*e.ReplacementYear)
		if m <= 0 {
			monthsText = fmt.Sprintf("%d ด.", m)
			monthsColor = "#E53935"
		} else {
			monthsText = fmt.Sprintf("%d ด.", m)
			if m <= 3 {
				monthsColor = "#E53935"
			} else if m <= 6 {
				monthsColor = "#FF9800"
			} else {
				monthsColor = "#2196F3"
			}
		}
	}

	bgColor := "#FFFFFF"
	if num%2 == 0 {
		bgColor = "#F5F5F5"
	}

	return map[string]interface{}{
		"type": "box", "layout": "vertical", "margin": "sm",
		"paddingAll": "8px", "backgroundColor": bgColor, "cornerRadius": "6px",
		"contents": []interface{}{
			map[string]interface{}{
				"type": "box", "layout": "horizontal",
				"contents": []interface{}{
					map[string]interface{}{
						"type": "text", "text": fmt.Sprintf("%d. %s", num, e.IDCode),
						"size": "xs", "weight": "bold", "flex": 4,
					},
					map[string]interface{}{
						"type": "text", "text": monthsText,
						"size": "xs", "color": monthsColor, "align": "end", "flex": 2, "weight": "bold",
					},
				},
			},
			map[string]interface{}{
				"type": "text", "text": fmt.Sprintf("   %s", modelName),
				"size": "xxs", "color": "#666666",
			},
		},
	}
}

// GetEquipmentExpiryByDeptFlex returns Flex Message for equipment expiry filtered by department
func GetEquipmentExpiryByDeptFlex(expired []entity.Equipment, nearExpiry []entity.Equipment, deptName string) map[string]interface{} {
	thisYear := time.Now().Year()
	nextYear := thisYear + 1

	bodyContents := []interface{}{}

	// Section 1: หมดอายุภายในปีนี้
	bodyContents = append(bodyContents, buildSectionHeader(
		fmt.Sprintf("🔴 หมดอายุภายในปี %d", thisYear), "#E53935"))
	bodyContents = append(bodyContents, buildEquipmentList(expired)...)
	bodyContents = append(bodyContents, map[string]interface{}{"type": "separator", "margin": "lg"})

	// Section 2: หมดอายุปีหน้า
	bodyContents = append(bodyContents, buildSectionHeader(
		fmt.Sprintf("🟡 หมดอายุปี %d", nextYear), "#FF9800"))
	bodyContents = append(bodyContents, buildEquipmentList(nearExpiry)...)

	return map[string]interface{}{
		"type": "bubble", "size": "mega",
		"header": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"backgroundColor": "#1B5E20", "paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text", "text": "📋 เครื่องใกล้หมดอายุ",
					"color": "#FFFFFF", "size": "lg", "weight": "bold",
				},
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("แผนก: %s", deptName),
					"color": "#FFFFFFCC", "size": "sm",
				},
				map[string]interface{}{
					"type": "text", "text": fmt.Sprintf("📅 ข้อมูล ณ %s", time.Now().Format("02/01/2006")),
					"color": "#FFFFFFCC", "size": "xs",
				},
			},
		},
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"spacing": "sm", "paddingAll": "12px",
			"contents": bodyContents,
		},
	}
}
