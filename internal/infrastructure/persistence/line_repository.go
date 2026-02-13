package persistence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"medical-webhook/internal/domain/line/model"
	"medical-webhook/internal/infrastructure/client"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// RepositoryImpl implements domain.repository.LineRepository
type LineRepository struct {
	client *client.Client
	token  string
}

// NewRepositoryImpl creates a new LINE repository implementation
func NewLineRepository(client *client.Client) *LineRepository {
	return &LineRepository{
		client: client,
		token:  client.Token,
	}
}

// ReplyMessage sends a reply message
func (r *LineRepository) ReplyMessage(replyToken, text string) error {
	if replyToken == "" {
		return nil
	}

	_, err := r.client.Bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{
				Text: text,
			},
		},
	})

	if err != nil {
		log.Printf("Error replying: %v", err)
		return err
	}

	return nil
}

// PushMessage sends a push message to a user
func (r *LineRepository) PushMessage(msg *model.OutgoingMessage) error {
	_, err := r.client.Bot.PushMessage(&messaging_api.PushMessageRequest{
		To: msg.To,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{
				Text: msg.Text,
			},
		},
	}, "")

	if err != nil {
		log.Printf("Error pushing message: %v", err)
		return err
	}

	return nil
}

// ReplyFlexMessage sends a Flex Message reply using raw HTTP API
func (r *LineRepository) ReplyFlexMessage(replyToken, altText string, flexContent map[string]interface{}) error {
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
func (r *LineRepository) PushFlexMessage(userID, altText string, flexContent map[string]interface{}) error {
	requestBody := map[string]interface{}{
		"to": userID,
		"messages": []map[string]interface{}{
			{"type": "flex", "altText": altText, "contents": flexContent},
		},
	}

	return r.sendRawJSON("https://api.line.me/v2/bot/message/push", requestBody)
}

// sendRawJSON sends a raw JSON request to LINE API
func (r *LineRepository) sendRawJSON(url string, body interface{}) error {
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

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("LINE API error %d: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("LINE API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Println("Flex message sent successfully")
	return nil
}

// BroadcastMessage - ส่งหาทุกคนที่เพิ่มเพื่อน Bot
func (r *LineRepository) BroadcastMessage(text string) error {
	requestBody := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"type": "text",
				"text": text,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/broadcast", bytes.NewBuffer(jsonBody))
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

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("LINE Broadcast API error %d: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("LINE Broadcast API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Println("Broadcast message sent successfully")
	return nil
}

// GetImageContent downloads image content from LINE using message ID
func (r *LineRepository) GetImageContent(messageID string) ([]byte, error) {
	log.Printf("Downloading image: %s", messageID)

	resp, err := r.client.BlobAPI.GetMessageContent(messageID)
	if err != nil {
		log.Printf("Error getting image content: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading image: %v", err)
		return nil, err
	}

	log.Printf("Downloaded image: %d bytes", len(imageBytes))
	return imageBytes, nil
}

func (r *LineRepository) GetProfile(userID string) (*model.UserProfile, error) {
	profile, err := r.client.Bot.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	return &model.UserProfile{
		UserID:        profile.UserId,
		DisplayName:   profile.DisplayName,
		PictureURL:    profile.PictureUrl,
		StatusMessage: profile.StatusMessage,
	}, nil
}

func (r *LineRepository) GetGroupMemberProfile(groupID, userID string) (*model.UserProfile, error) {
	profile, err := r.client.Bot.GetGroupMemberProfile(groupID, userID)
	if err != nil {
		return nil, err
	}

	return &model.UserProfile{
		UserID:        profile.UserId,
		DisplayName:   profile.DisplayName,
		PictureURL:    profile.PictureUrl,
		StatusMessage: "", // GroupMemberProfileResponse might not have status message
	}, nil
}

func (r *LineRepository) GetRoomMemberProfile(roomID, userID string) (*model.UserProfile, error) {
	profile, err := r.client.Bot.GetRoomMemberProfile(roomID, userID)
	if err != nil {
		return nil, err
	}

	return &model.UserProfile{
		UserID:        profile.UserId,
		DisplayName:   profile.DisplayName,
		PictureURL:    profile.PictureUrl,
		StatusMessage: "", // RoomMemberProfileResponse might not have status message
	}, nil
}
