package model

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeLocation MessageType = "location"
)

// IncomingMessage represents an incoming LINE message
type IncomingMessage struct {
	UserID      string
	GroupID     string
	SourceType  string
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

// OutgoingMessage represents an outgoing message
type OutgoingMessage struct {
	To   string
	Text string
}
