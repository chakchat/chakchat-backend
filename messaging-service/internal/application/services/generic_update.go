package services

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/google/uuid"
)

const (
	UpdateTypeTextMessage       = "text_message"
	UpdateTypeTextMessageEdited = "text_message_edited"
	UpdateTypeFileMessage       = "file_message"
	UpdateTypeReaction          = "reaction"
	UpdateTypeDeleted           = "update_deleted"
)

type GenericUpdate struct {
	UpdateID   int64
	ChatID     uuid.UUID
	SenderID   uuid.UUID
	UpdateType string

	CreatedAt int64
	Info      GenericUpdateInfo
}

// You should get one of fields depending on GenericUpdate.UpdateType
type GenericUpdateInfo struct {
	TextMessage       *TextMessageInfo
	TextMessageEdited *TextMessageEditedInfo
	FileMessage       *FileMessageInfo
	Deleted           *DeletedInfo
	Reaction          *ReactionInfo
}

type TextMessageInfo struct {
	Text    string
	Edited  *GenericUpdate
	ReplyTo *int64
}

type TextMessageEditedInfo struct {
	MessageID int64
	NewText   string
}

type FileMessageInfo struct {
	File    dto.FileMetaDTO
	ReplyTo *int64
}

type DeletedInfo struct {
	DeletedID  int64
	DeleteMode string
}

type ReactionInfo struct {
	Reaction string
}
