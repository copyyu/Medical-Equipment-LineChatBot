package line

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"medical-webhook/internal/domain/line/model"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// RepositoryImpl implements domain.repository.LineRepository
type RepositoryImpl struct {
	client *Client
	token  string
}

// NewRepositoryImpl creates a new LINE repository implementation
func NewRepositoryImpl(client *Client) *RepositoryImpl {
	return &RepositoryImpl{
		client: client,
		token:  client.token,
	}
}

// ReplyMessage sends a reply message
func (r *RepositoryImpl) ReplyMessage(replyToken, text string) error {
	if replyToken == "" {
		return nil
	}

	_, err := r.client.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{
				Text: text,
			},
		},
	})

	if err != nil {
		log.Printf("❌ Error replying: %v", err)
		return err
	}

	return nil
}

// PushMessage sends a push message to a user
func (r *RepositoryImpl) PushMessage(msg *model.OutgoingMessage) error {
	_, err := r.client.bot.PushMessage(&messaging_api.PushMessageRequest{
		To: msg.To,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{
				Text: msg.Text,
			},
		},
	}, "")

	if err != nil {
		log.Printf("❌ Error pushing message: %v", err)
		return err
	}

	return nil
}

// ReplyFlexMessage sends a Flex Message reply using raw HTTP API
func (r *RepositoryImpl) ReplyFlexMessage(replyToken, altText string, flexContent map[string]interface{}) error {
	if replyToken == "" {
		return nil
	}

	requestBody := map[string]interface{}{
		"replyToken": replyToken,
		"messages": []map[string]interface{}{
			{"type": "flex", "altText": altText, "contents": flexContent},
		},
	}

	return r.sendRawJSON("https://api.line.me/v2/bot/message/reply", requestBody)
}

// PushFlexMessage sends a Flex Message push to a user
func (r *RepositoryImpl) PushFlexMessage(userID, altText string, flexContent map[string]interface{}) error {
	requestBody := map[string]interface{}{
		"to": userID,
		"messages": []map[string]interface{}{
			{"type": "flex", "altText": altText, "contents": flexContent},
		},
	}

	return r.sendRawJSON("https://api.line.me/v2/bot/message/push", requestBody)
}

// sendRawJSON sends a raw JSON request to LINE API
func (r *RepositoryImpl) sendRawJSON(url string, body interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("❌ LINE API error %d: %s", resp.StatusCode, string(bodyBytes))
	} else {
		log.Println("✅ Flex message sent successfully")
	}

	return nil
}
