package repository

import "medical-webhook/internal/domain/line/model"

// LineRepository defines interface for LINE messaging operations
type LineRepository interface {
	ReplyMessage(replyToken, text string) error
	PushMessage(msg *model.OutgoingMessage) error
	ReplyFlexMessage(replyToken, altText string, flexContent map[string]interface{}) error
	PushFlexMessage(userID, altText string, flexContent map[string]interface{}) error
	BroadcastMessage(text string) error
	BroadcastFlexMessage(altText string, flexContent map[string]interface{}) error
	GetImageContent(messageID string) ([]byte, error)
	GetProfile(userID string) (*model.UserProfile, error)
	GetGroupMemberProfile(groupID, userID string) (*model.UserProfile, error)
	GetRoomMemberProfile(roomID, userID string) (*model.UserProfile, error)
}
