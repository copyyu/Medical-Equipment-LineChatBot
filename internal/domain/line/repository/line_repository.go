package repository

import "medical-webhook/internal/domain/line/model"

// LineRepository defines interface for LINE messaging operations
type LineRepository interface {
	ReplyMessage(replyToken, text string) error
	PushMessage(msg *model.OutgoingMessage) error
}
