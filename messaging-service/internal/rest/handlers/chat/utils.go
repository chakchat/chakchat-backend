package chat

import (
	"fmt"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/google/uuid"
)

func convertGenericChatResp(chat *services.GenericChat) any {
	resp := struct {
		ChatID    uuid.UUID   `json:"chat_id"`
		Type      string      `json:"type"`
		Members   []uuid.UUID `json:"members"`
		CreatedAt int64       `json:"created_at"`
		Info      any         `json:"info"`
	}{
		ChatID:    chat.ChatID,
		Type:      chat.ChatType,
		Members:   chat.Members,
		CreatedAt: chat.CreatedAt,
	}

	switch chat.ChatType {
	case services.ChatTypePersonal:
		resp.Info = struct {
			BlockedBy []uuid.UUID `json:"blocked_by"`
		}{
			BlockedBy: chat.Info.Personal().BlockedBy,
		}
	case services.ChatTypeGroup:
		resp.Info = struct {
			Admin       uuid.UUID `json:"admin_id"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			GroupPhoto  string    `json:"group_photo"`
		}{
			Admin:       chat.Info.Group().Admin,
			Name:        chat.Info.Group().GroupName,
			Description: chat.Info.Group().GroupDescription,
			GroupPhoto:  chat.Info.Group().GroupPhoto,
		}
	case services.ChatTypeSecretPersonal:
		resp.Info = struct {
			Expiration *time.Duration `json:"expiration"`
		}{
			Expiration: chat.Info.SecretPersonal().Expiration,
		}
	case services.ChatTypeSecretGroup:
		resp.Info = struct {
			Admin       uuid.UUID      `json:"admin_id"`
			Name        string         `json:"name"`
			Description string         `json:"description"`
			GroupPhoto  string         `json:"group_photo"`
			Expiration  *time.Duration `json:"expiration"`
		}{
			Admin:       chat.Info.SecretGroup().Admin,
			Name:        chat.Info.SecretGroup().GroupName,
			Description: chat.Info.SecretGroup().GroupDescription,
			GroupPhoto:  chat.Info.SecretGroup().GroupPhoto,
			Expiration:  chat.Info.SecretPersonal().Expiration,
		}
	default:
		panic(fmt.Errorf("uknown chat type: %s", chat.ChatType))
	}

	return resp
}
