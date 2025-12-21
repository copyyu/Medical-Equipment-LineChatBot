package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// NotifyClient handles LINE Notify API operations
type NotifyClient struct {
	Token string
}

// NewNotifyClient creates a new LINE Notify client
func NewNotifyClient(token string) *NotifyClient {
	return &NotifyClient{Token: token}
}

// SendMessage sends a message via LINE Notify
func (c *NotifyClient) SendMessage(message string) error {
	apiURL := "https://notify-api.line.me/api/notify"

	data := url.Values{}
	data.Set("message", message)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("LINE Notify API error %d: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("LINE Notify API returned status: %d", resp.StatusCode)
	}

	log.Println("LINE Notify sent successfully")
	return nil
}

// FormatReplacementAlert formats equipment replacement alert message
func FormatReplacementAlert(alerts []interface{}, notifyRound string) string {
	roundName := "มิถุนายน (6 เดือนก่อน)"
	if notifyRound == "AUGUST" {
		roundName = "สิงหาคม (1 ปีก่อน)"
	}

	message := fmt.Sprintf("\n🔔 แจ้งเตือนอุปกรณ์ที่ใกล้ครบกำหนดเปลี่ยน\n")
	message += fmt.Sprintf("📅 รอบแจ้งเตือน: %s\n", roundName)
	message += "━━━━━━━━━━━━━━━━━━━━━━\n\n"

	if len(alerts) == 0 {
		message += "✅ ไม่มีอุปกรณ์ที่ต้องแจ้งเตือนในรอบนี้"
		return message
	}

	for i, item := range alerts {
		alert := item.(map[string]interface{})

		message += fmt.Sprintf("%d. 📦 %s\n", i+1, alert["id_code"].(string))
		message += fmt.Sprintf("   %s - %s\n", alert["brand_name"].(string), alert["model_name"].(string))
		message += fmt.Sprintf("   แผนก: %s\n", alert["department_name"].(string))
		message += fmt.Sprintf("   ⏰ เหลืออีก %d เดือน\n", int(alert["months_remaining"].(float64)))

		if i < len(alerts)-1 {
			message += "\n"
		}
	}

	message += fmt.Sprintf("\n━━━━━━━━━━━━━━━━━━━━━━")
	message += fmt.Sprintf("\n📊 รวมทั้งหมด: %d รายการ", len(alerts))

	return message
}
