package line

import (
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// Client wraps LINE Messaging API client
type Client struct {
	bot *messaging_api.MessagingApiAPI
}

// NewClient creates a new LINE client
func NewClient(channelToken string) (*Client, error) {
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return nil, err
	}

	return &Client{
		bot: bot,
	}, nil
}

// GetBot returns the underlying LINE bot instance
func (c *Client) GetBot() *messaging_api.MessagingApiAPI {
	return c.bot
}
