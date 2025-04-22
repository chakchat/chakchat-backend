// generic package presents unified format of generic updates and chats
package generic

import (
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type Chat struct {
	ChatID    uuid.UUID `json:"chat_id"`
	CreatedAt int64 	`json:"created_at"`
	Type  string `json:"type"`
	Members   []uuid.UUID `json:"members"`
	Info      ChatInfo `json:"info"`
	// Last update ID in the chat.
	// Be careful, it may hold even ID of update not visible for user (e.g. deleted)
	LastUpdateID *int64 `json:"last_update_id,omitempty"`
	// Holds last updates to show chat preview in the client.
	// Not fetched by default.
	UpdatePreview []Update `json:"update_preview,omitempty"`
}

type ChatInfo struct {
	Personal *PersonalInfo `json:",inline,omitempty"`
	Group *GroupInfo `json:",inline,omitempty"`
	SecretPersonal *SecretPersonalInfo `json:",inline,omitempty"`
	SecretGroup *SecretGroupInfo `json:",inline,omitempty"`
}

type PersonalInfo struct {
	BlockedBy []uuid.UUID `json:"blocked_by"`
}

type GroupInfo struct {
	AdminID            uuid.UUID `json:"admin_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupPhoto       string `json:"group_photo"`
}

type SecretPersonalInfo struct {
	Expiration *time.Duration `json:"expiration"`
}

type SecretGroupInfo struct {
	AdminID            uuid.UUID `json:"admin_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupPhoto       string `json:"group_photo"`
	Expiration       *time.Duration `json:"expiration"`
}

func FromPersonalChatDTO(chatDTO *dto.PersonalChatDTO) Chat {
	return Chat{
		ChatID:      chatDTO.ID,
		CreatedAt:   chatDTO.CreatedAt,
		Type:        domain.ChatTypePersonal,
		Members:     chatDTO.Members[:],
		Info: ChatInfo{
			Personal: &PersonalInfo{
				BlockedBy: chatDTO.BlockedBy,
			},
		},
		LastUpdateID:  nil,
		UpdatePreview: nil,
	}
}

func FromGroupChatDTO(chatDTO *dto.GroupChatDTO) Chat {
	return Chat{
		ChatID:    chatDTO.ID,
		CreatedAt: chatDTO.CreatedAt,
		Type:      domain.ChatTypeGroup,
		Members:   chatDTO.Members,
		Info: ChatInfo{
			Group: &GroupInfo{
				AdminID:     chatDTO.Admin,
				Name:        chatDTO.Name,
				Description: chatDTO.Description,
				GroupPhoto:  chatDTO.GroupPhoto,
			},
		},
		LastUpdateID:  nil,
		UpdatePreview: nil,
	}
}

func FromSecretPersonalChatDTO(chatDTO *dto.SecretPersonalChatDTO) Chat {
	return Chat{
		ChatID:    chatDTO.ID,
		CreatedAt: chatDTO.CreatedAt,
		Type:      domain.ChatTypeSecretPersonal,
		Members:   chatDTO.Members[:],
		Info: ChatInfo{
			SecretPersonal: &SecretPersonalInfo{
				Expiration: chatDTO.Expiration,
			},
		},
		LastUpdateID:  nil,
		UpdatePreview: nil,
	}
}

func FromSecretGroupChatDTO(chatDTO *dto.SecretGroupChatDTO) Chat {
	return Chat{
		ChatID:    chatDTO.ID,
		CreatedAt: chatDTO.CreatedAt,
		Type:      domain.ChatTypeSecretGroup,
		Members:   chatDTO.Members,
		Info: ChatInfo{
			SecretGroup: &SecretGroupInfo{
				AdminID:     chatDTO.Admin,
				Name:        chatDTO.Name,
				Description: chatDTO.Description,
				GroupPhoto:  chatDTO.GroupPhotoURL,
				Expiration:  chatDTO.Expiration,
			},
		},
		LastUpdateID:  nil,
		UpdatePreview: nil,
	}
}
