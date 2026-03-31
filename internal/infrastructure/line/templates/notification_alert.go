package templates

// GetNotificationAlertFlex returns a Flex message containing the notification text and a download button
func GetNotificationAlertFlex(notificationText string, downloadURL string) map[string]interface{} {
	return map[string]interface{}{
		"type": "bubble",
		"size": "mega",
		"body": map[string]interface{}{
			"type":       "box",
			"layout":     "vertical",
			"paddingAll": "15px",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   notificationText,
					"wrap":   true,
					"size":   "sm",
					"color":  ColorText,
					"weight": "regular",
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
					"color": ColorSuccess,
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "📥 ดาวน์โหลด Excel",
						"uri":   downloadURL,
					},
				},
			},
		},
	}
}
