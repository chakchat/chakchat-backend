package generic

import (
	"encoding/base64"
	"encoding/json"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type Update struct {
	UpdateID   int64     `json:"update_id"`
	ChatID     uuid.UUID `json:"chat_id"`
	SenderID   uuid.UUID `json:"sender_id"`
	UpdateType string    `json:"type"`

	CreatedAt int64         `json:"created_at"`
	Content   UpdateContent `json:"content"`
}

// You should get one of fields depending on GenericUpdate.UpdateType
type UpdateContent struct {
	TextMessage       *TextMessageContent
	TextMessageEdited *TextMessageEditedContent
	FileMessage       *FileMessageContent
	Deleted           *DeletedContent
	Reaction          *ReactionContent 
	Secret            *SecretUpdateContent
}

func (c UpdateContent) MarshalJSON() ([]byte, error) {
	switch {
	case c.TextMessage != nil:
		return json.Marshal(c.TextMessage)
	case c.TextMessageEdited != nil:
		return json.Marshal(c.TextMessageEdited)
	case c.FileMessage != nil:
		return json.Marshal(c.FileMessage)
	case c.Deleted != nil:
		return json.Marshal(c.Deleted)
	case c.Reaction != nil:
		return json.Marshal(c.Reaction)
	case c.Secret != nil:
		return json.Marshal(c.Secret)
	default:
		return nil, nil
	}
}

type TextMessageContent struct {
	Text      string   `json:"text"`
	Edited    *Update  `json:"edited,omitempty"`
	ReplyTo   *int64   `json:"reply_to,omitempty"`
	Reactions []Update `json:"reactions,omitempty"`
}

type TextMessageEditedContent struct {
	MessageID int64  `json:"message_id"`
	NewText   string `json:"new_text"`
}

type FileMessageContent struct {
	File      FileMeta `json:"file"`
	ReplyTo   *int64   `json:"reply_to,omitempty"`
	Reactions []Update `json:"reactions,omitempty"`
}

type FileMeta struct {
	FileId    uuid.UUID `json:"file_id"`
	FileName  string    `json:"file_name"`
	MimeType  string    `json:"mime_type"`
	FileSize  int64     `json:"file_size"`
	FileURL   string    `json:"file_url"`
	CreatedAt int64     `json:"created_at"`
}

type DeletedContent struct {
	DeletedID   int64  `json:"deleted_id"`
	DeletedMode string `json:"deleted_mode"`
}

type ReactionContent struct {
	Reaction  string `json:"reaction"`
	MessageID int64  `json:"message_id"`
}

type SecretUpdateContent struct {
	PayloadBase64              string `json:"payload"`
	InitializationVectorBase64 string `json:"initialization_vector"`
	KeyHashBase64              string `json:"key_hash"`
}

func FromFileMessageDTO(msg *dto.FileMessageDTO) Update {
	var replyTo *int64
	if msg.ReplyTo != nil {
		cp := int64(*msg.ReplyTo)
		replyTo = &cp
	}

	return Update{
		UpdateID:   msg.UpdateID,
		ChatID:     msg.ChatID,
		SenderID:   msg.SenderID,
		UpdateType: domain.UpdateTypeFileMessage,
		CreatedAt:  msg.CreatedAt,
		Content: UpdateContent{
			FileMessage: &FileMessageContent{
				File: FileMeta{
					FileId:    msg.File.FileId,
					FileName:  msg.File.FileName,
					MimeType:  msg.File.MimeType,
					FileSize:  msg.File.FileSize,
					FileURL:   msg.File.FileURL,
					CreatedAt: msg.File.CreatedAt,
				},
				ReplyTo:   replyTo,
				Reactions: nil,
			},
		},
	}
}

func FromTextMessageDTO(msg *dto.TextMessageDTO) Update {
	var replyTo *int64
	if msg.ReplyTo != nil {
		cp := int64(*msg.ReplyTo)
		replyTo = &cp
	}

	var edited *Update
	if msg.Edited != nil {
		editedUpdate := FromTextMessageEditedDTO(msg.Edited)
		edited = &editedUpdate
	}

	return Update{
		UpdateID:   msg.UpdateID,
		ChatID:     msg.ChatID,
		SenderID:   msg.SenderID,
		UpdateType: domain.UpdateTypeTextMessage,
		CreatedAt:  msg.CreatedAt,
		Content: UpdateContent{
			TextMessage: &TextMessageContent{
				Text:      msg.Text,
				Edited:    edited,
				ReplyTo:   replyTo,
				Reactions: nil,
			},
		},
	}
}

func FromTextMessageEditedDTO(msg *dto.TextMessageEditedDTO) Update {
	return Update{
		UpdateID:   msg.UpdateID,
		ChatID:     msg.ChatID,
		SenderID:   msg.SenderID,
		UpdateType: domain.UpdateTypeTextMessageEdited,
		CreatedAt:  msg.CreatedAt,
		Content: UpdateContent{
			TextMessageEdited: &TextMessageEditedContent{
				MessageID: msg.MessageID,
				NewText:   msg.NewText,
			},
		},
	}
}

func FromUpdateDeletedDTO(msg *dto.UpdateDeletedDTO) Update {
	return Update{
		UpdateID:   msg.UpdateID,
		ChatID:     msg.ChatID,
		SenderID:   msg.SenderID,
		UpdateType: domain.UpdateTypeDeleted,
		CreatedAt:  msg.CreatedAt,
		Content: UpdateContent{
			Deleted: &DeletedContent{
				DeletedID:   msg.DeletedID,
				DeletedMode: msg.DeleteMode,
			},
		},
	}
}

func FromSecretUpdateDTO(msg *dto.SecretUpdateDTO) Update {
	return Update{
		UpdateID:   msg.UpdateID,
		ChatID:     msg.ChatID,
		SenderID:   msg.SenderID,
		UpdateType: domain.UpdateTypeSecret,
		CreatedAt:  msg.CreatedAt,
		Content: UpdateContent{
			Secret: &SecretUpdateContent{
				PayloadBase64:              base64.StdEncoding.EncodeToString(msg.Payload),
				InitializationVectorBase64: base64.StdEncoding.EncodeToString(msg.InitializationVector),
				KeyHashBase64:              base64.StdEncoding.EncodeToString(msg.KeyHash),
			},
		},
	}
}

func FromReactionDTO(r *dto.ReactionDTO) Update {
	return Update{
		UpdateID:   r.UpdateID,
		ChatID:     r.ChatID,
		SenderID:   r.SenderID,
		UpdateType: domain.UpdateTypeReaction,
		CreatedAt:  r.CreatedAt,
		Content: UpdateContent{
			Reaction: &ReactionContent{
				Reaction:  r.ReactionType,
				MessageID: r.MessageID,
			},
		},
	}
}
