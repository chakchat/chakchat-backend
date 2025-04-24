package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type Notification struct {
	Receivers uuid.UUID       `json:"receivers"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
}

type UpdateMessage struct {
	ChatId    uuid.UUID       `json:"chat_id"`
	UpdateID  uuid.UUID       `json:"update_id"`
	Type      string          `json:"type"`
	SenderID  uuid.UUID       `json:"sender_id"`
	CreatedAt time.Time       `json:"created_at"`
	Content   json.RawMessage `json:"content"`
}

type TextMessageContent struct {
	Text    string    `json:"text"`
	ReplyTo uuid.UUID `json:"reply_to"`
}

type FileMessageContent struct {
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
}

type ReactionMessageContent struct {
	Reaction string `json:"reaction"`
}

type DeleteMessageContent struct {
	DeletedMode string `json:"deleted_mode"`
}

type CreateChatMessage struct {
	SenderID uuid.UUID `json:"sender_id"`
	Chat     *struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}
}

type GroupInfoUpdated struct {
	SenderID    uuid.UUID `json:"sender_id"`
	ChatID      uuid.UUID `json:"chat_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupPhoto  string    `json:"group_photo"`
}

type UpdateGroupMembers struct {
	SenderID uuid.UUID   `json:"sender_id"`
	ChatID   uuid.UUID   `json:"chat_id"`
	Members  []uuid.UUID `json:"members"`
}

type GRPCClients interface {
	GetChatType() (string, error)
	GetGroupName() (string, error)
	GetName(ctx context.Context, userId uuid.UUID) (*string, error)
}

type Parser struct {
	grpcHandler GRPCClients
}

func NewParser(grpcHandl GRPCClients) *Parser {
	return &Parser{
		grpcHandler: grpcHandl,
	}
}

func (p *Parser) ParseNotification(ctx context.Context, raw []byte) (string, error) {
	var notific Notification
	if err := json.Unmarshal(raw, &notific); err != nil {
		return "", err
	}

	switch notific.Type {
	case "update":
		return p.ParseUpdateNotification(ctx, notific.Data)
	case "chat_created":
		return p.ParseChatCreated(ctx, notific.Data)
	case "group_info_updated":
		return p.ParseGroupInfoUpdated(ctx, notific.Data)
	case "group_members_added", "group_members_removed":
		return p.ParseGroupMembersChanged(ctx, notific.Type, notific.Data)
	}
	return "", nil
}

func (p *Parser) ParseUpdateNotification(ctx context.Context, data json.RawMessage) (string, error) {
	var update UpdateMessage
	if err := json.Unmarshal(data, &update); err != nil {
		return "", err
	}
	switch update.Type {
	case "text_message":
		var content TextMessageContent
		if err := json.Unmarshal(update.Content, &content); err != nil {
			return "", nil
		}

		chatType, err := p.grpcHandler.GetChatType()
		if err != nil {
			return "", err
		}
		sender, err := p.grpcHandler.GetName(ctx, update.SenderID)
		if err != nil {
			return "", err
		}
		if chatType == "group" {
			groupName, err := p.grpcHandler.GetGroupName()
			if err != nil {
				return "", nil
			}
			return fmt.Sprintf("%s sent new message: %s from %s", sender, Truncate(content.Text, 30), groupName), nil
		}
		return fmt.Sprintf("%s sent new message: %s", sender, Truncate(content.Text, 30)), nil
	case "file":
		var content FileMessageContent
		if err := json.Unmarshal(update.Content, &content); err != nil {
			return "", err
		}
		chatType, err := p.grpcHandler.GetChatType()
		if err != nil {
			return "", err
		}
		sender, err := p.grpcHandler.GetName(ctx, update.SenderID)
		if err != nil {
			return "", err
		}
		if chatType == "group" {
			groupName, err := p.grpcHandler.GetGroupName()
			if err != nil {
				return "", nil
			}
			return fmt.Sprintf("%s sent new filr: %s from %s", sender, content.FileName, groupName), nil
		}
		return fmt.Sprintf("%s sent new file: %s", sender, content.FileName), nil
	case "reaction":
		var content ReactionMessageContent
		if err := json.Unmarshal(update.Content, &content); err != nil {
			return "", err
		}
		sender, err := p.grpcHandler.GetName(ctx, update.SenderID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s put new reaction: %s", sender, content.Reaction), nil
	case "delete":
		return "", nil
	}
	return "", fmt.Errorf("incorrect json")
}

func (p *Parser) ParseChatCreated(ctx context.Context, data json.RawMessage) (string, error) {
	var chat CreateChatMessage
	if err := json.Unmarshal(data, &chat); err != nil {
		return "", err
	}
	return fmt.Sprintf("New %s chat: %s", chat.Chat.Type, chat.Chat.Name), nil
}

func (p *Parser) ParseGroupInfoUpdated(ctx context.Context, data json.RawMessage) (string, error) {
	var groupInfo GroupInfoUpdated
	if err := json.Unmarshal(data, &groupInfo); err != nil {
		return "", nil
	}
	sender, err := p.grpcHandler.GetName(ctx, groupInfo.SenderID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s changed group info in %s", sender, groupInfo.Name), nil
}

func (p *Parser) ParseGroupMembersChanged(ctx context.Context, notifiqType string, data json.RawMessage) (string, error) {
	var group UpdateGroupMembers
	if err := json.Unmarshal(data, &group); err != nil {
		return "", nil
	}

	groupName, err := p.grpcHandler.GetGroupName()
	if err != nil {
		return "", nil
	}
	sender, err := p.grpcHandler.GetName(ctx, group.SenderID)
	if err != nil {
		return "", err
	}
	switch notifiqType {
	case "group_members_added":
		return fmt.Sprintf("%s added new members in %s", sender, groupName), nil
	case "group_members_removed":
		return fmt.Sprintf("%s removed new members in %s", sender, groupName), nil
	}
	return "", nil
}

func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
