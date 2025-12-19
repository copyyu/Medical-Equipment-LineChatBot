package model

// MessageType represents type of message
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeLocation MessageType = "location"
)

// IncomingMessage represents message from LINE
type IncomingMessage struct {
	UserID      string
	MessageType MessageType
	Text        string
	ImageID     string
	Location    *Location
	ReplyToken  string
}

// Location represents location data
type Location struct {
	Latitude  float64
	Longitude float64
	Address   string
}

// OutgoingMessage represents message to send
type OutgoingMessage struct {
	To   string
	Text string
}
