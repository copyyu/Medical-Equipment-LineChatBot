package client

import (
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// Client wraps the LINE Messaging API client
type Client struct {
	Bot     *messaging_api.MessagingApiAPI
	BlobAPI *messaging_api.MessagingApiBlobAPI
	Token   string
}

// NewClient creates a new LINE client
func NewClient(channelToken string) (*Client, error) {
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return nil, err
	}

	blobAPI, err := messaging_api.NewMessagingApiBlobAPI(channelToken)
	if err != nil {
		return nil, err
	}

	return &Client{
		Bot:     bot,
		BlobAPI: blobAPI,
		Token:   channelToken,
	}, nil
}
