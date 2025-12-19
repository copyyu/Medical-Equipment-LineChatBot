package line

import (
	"log"
	"medical-webhook/internal/domain/line/model"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// RepositoryImpl implements domain.repository.LineRepository
type RepositoryImpl struct {
	client *Client
}

// NewRepositoryImpl creates a new LINE repository implementation
func NewRepositoryImpl(client *Client) *RepositoryImpl {
	return &RepositoryImpl{
		client: client,
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
