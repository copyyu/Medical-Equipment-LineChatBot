package line

import (
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// Client wraps the LINE Messaging API client
type Client struct {
	bot   *messaging_api.MessagingApiAPI
	token string
}

// NewClient creates a new LINE client
func NewClient(channelToken string) (*Client, error) {
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return nil, err
	}

	return &Client{
		bot:   bot,
		token: channelToken,
	}, nil
}
