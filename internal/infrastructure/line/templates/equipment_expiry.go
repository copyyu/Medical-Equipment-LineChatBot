package templates

import (
	"fmt"
	"medical-webhook/internal/domain/line/entity"
	"time"
)

const lineDisplayLimit = 10 // จำนวน item สูงสุดที่แสดงใน LINE Flex

// GetEquipmentExpiryFlex - แสดง Flex รายการเครื่องใกล้หมดอายุ (ทุกแผนก)
func GetEquipmentExpiryFlex(expired []entity.Equipment, nearExpiry []entity.Equipment, baseURL string) map[string]interface{} {
	thisYear := time.Now().Year()
	nextYear := thisYear + 1

	totalCount := len(expired) + len(nearExpiry)
	downloadURL := baseURL + "/notifications/export/expiry"

	bodyContents := []interface{}{}
	bodyContents = append(bodyContents, buildSectionHeader(
		fmt.Sprintf("🔴 หมดอายุภายในปี %d", thisYear), "#E53935"))
	bodyContents = append(bodyContents, buildEquipmentList(expired, lineDisplayLimit)...)
	bodyContents = append(bodyContents, map[string]interface{}{"type": "separator", "margin": "lg"})
	bodyContents = append(bodyContents, buildSectionHeader(
		fmt.Sprintf("🟡 หมดอายุปี %d", nextYear), "#FF9800"))
	bodyContents = append(bodyContents, buildEquipmentList(nearExpiry, lineDisplayLimit)...)

	return map[string]interface{}{
		"type": "bubble", "size": "mega",
		"header": buildExpiryHeader("📋 รายการเครื่องใกล้หมดอายุ", "", totalCount, time.Now()),
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"spacing": "sm", "paddingAll": "12px",
			"contents": bodyContents,
		},
		"footer": buildExportFooter(downloadURL, totalCount),
	}
}

// GetEquipmentExpiryByDeptFlex - แสดง Flex รายการเครื่องใกล้หมดอายุของแผนก
func GetEquipmentExpiryByDeptFlex(expired []entity.Equipment, nearExpiry []entity.Equipment, deptName string, departmentID uint, baseURL string) map[string]interface{} {
	thisYear := time.Now().Year()
	nextYear := thisYear + 1

	totalCount := len(expired) + len(nearExpiry)
	downloadURL := fmt.Sprintf("%s/notifications/export/expiry?dept_id=%d", baseURL, departmentID)

	bodyContents := []interface{}{}
	bodyContents = append(bodyContents, buildSectionHeader(
		fmt.Sprintf("🔴 หมดอายุภายในปี %d", thisYear), "#E53935"))
	bodyContents = append(bodyContents, buildEquipmentList(expired, lineDisplayLimit)...)
	bodyContents = append(bodyContents, map[string]interface{}{"type": "separator", "margin": "lg"})
	bodyContents = append(bodyContents, buildSectionHeader(
		fmt.Sprintf("🟡 หมดอายุปี %d", nextYear), "#FF9800"))
	bodyContents = append(bodyContents, buildEquipmentList(nearExpiry, lineDisplayLimit)...)

	return map[string]interface{}{
		"type": "bubble", "size": "mega",
		"header": buildExpiryHeader("📋 เครื่องใกล้หมดอายุ", deptName, totalCount, time.Now()),
		"body": map[string]interface{}{
			"type": "box", "layout": "vertical",
			"spacing": "sm", "paddingAll": "12px",
			"contents": bodyContents,
		},
		"footer": buildExportFooter(downloadURL, totalCount),
	}
}

// buildExpiryHeader สร้าง header ของ Flex Message แสดงชื่อ แผนก และ count รวม
func buildExpiryHeader(title string, deptName string, totalCount int, now time.Time) map[string]interface{} {
	contents := []interface{}{
		map[string]interface{}{
			"type": "text", "text": title,
			"color": "#FFFFFF", "size": "lg", "weight": "bold",
		},
	}
	if deptName != "" {
		contents = append(contents, map[string]interface{}{
			"type": "text", "text": fmt.Sprintf("แผนก: %s", deptName),
			"color": "#FFFFFFCC", "size": "sm",
		})
	}
	displayedNote := ""
	if totalCount > lineDisplayLimit*2 {
		displayedNote = fmt.Sprintf("แสดง %d จาก %d รายการ (Export Excel เพื่อดูทั้งหมด)", lineDisplayLimit*2, totalCount)
	} else {
		displayedNote = fmt.Sprintf("รวม %d รายการ | %s", totalCount, now.Format("02/01/2006"))
	}
	contents = append(contents, map[string]interface{}{
		"type": "text", "text": displayedNote,
		"color": "#FFFFFFCC", "size": "xxs", "wrap": true,
	})

	return map[string]interface{}{
		"type": "box", "layout": "vertical",
		"backgroundColor": "#1B5E20", "paddingAll": "15px",
		"contents": contents,
	}
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

// buildEquipmentList สร้างรายการเครื่องมือ
func buildEquipmentList(equipments []entity.Equipment, displayLimit int) []interface{} {
	items := []interface{}{}

	if len(equipments) == 0 {
		items = append(items, map[string]interface{}{
			"type": "text", "text": "✅ ไม่มีเครื่องในหมวดนี้",
			"size": "xs", "color": "#4CAF50", "align": "center", "margin": "sm",
		})
		return items
	}

	displayItems := equipments
	hasMore := false
	remaining := 0
	if len(equipments) > displayLimit {
		displayItems = equipments[:displayLimit]
		hasMore = true
		remaining = len(equipments) - displayLimit
	}

	items = append(items, map[string]interface{}{
		"type": "text", "text": fmt.Sprintf("📊 รวม %d รายการ", len(equipments)),
		"size": "xs", "color": "#888888", "margin": "sm",
	})

	for i, e := range displayItems {
		items = append(items, buildEquipmentRow(i+1, e))
	}

	if hasMore {
		items = append(items, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("...และอีก %d รายการ → กด Export Excel เพื่อดูทั้งหมด", remaining),
			"size": "xxs", "color": "#E53935", "align": "center",
			"margin": "sm", "wrap": true,
		})
	}

	return items
}

// buildEquipmentRow แต่ละรายการ: running number + ID code + ชื่อเครื่อง + เหลือกี่เดือน
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
			switch {
			case m <= 3:
				monthsColor = "#E53935"
			case m <= 6:
				monthsColor = "#FF9800"
			default:
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

func calcMonthsRemaining(replacementYear int) int {
	now := time.Now()
	return (replacementYear-now.Year())*12 - int(now.Month()) + 8
}

func buildExportFooter(downloadURL string, totalCount int) map[string]interface{} {
	label := "Export Excel"
	return map[string]interface{}{
		"type": "box", "layout": "vertical",
		"paddingAll": "12px",
		"contents": []interface{}{
			map[string]interface{}{
				"type":   "button",
				"style":  "primary",
				"color":  "#1B5E20",
				"height": "sm",
				"action": map[string]interface{}{
					"type":  "uri",
					"label": label,
					"uri":   downloadURL,
				},
			},
		},
	}
}
