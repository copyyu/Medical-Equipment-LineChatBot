package client

import (
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// Client wraps the LINE Messaging API client
type Client struct {
	Bot     *messaging_api.MessagingApiAPI
	BlobAPI *messaging_api.MessagingApiBlobAPI
	Token   string
}

// NewClient creates a new LINE client. A timeout is applied to the underlying
// HTTP client so a slow/stuck LINE API call can never hang a request goroutine
// indefinitely (the SDK's default client has no timeout).
func NewClient(channelToken string, timeout time.Duration) (*Client, error) {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	httpClient := &http.Client{Timeout: timeout}

	bot, err := messaging_api.NewMessagingApiAPI(channelToken, messaging_api.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	blobAPI, err := messaging_api.NewMessagingApiBlobAPI(channelToken, messaging_api.WithBlobHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &Client{
		Bot:     bot,
		BlobAPI: blobAPI,
		Token:   channelToken,
	}, nil
}
