package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/group"
	"github.com/google/uuid"
)

type GroupChatDTO struct {
	ID      uuid.UUID
	Admin   uuid.UUID
	Members []uuid.UUID

	Name        string
	Description string
	GroupPhoto  string
	CreatedAt   int64
}

func NewGroupChatDTO(g *group.GroupChat) GroupChatDTO {
	return GroupChatDTO{
		ID:          uuid.UUID(g.ID),
		Admin:       uuid.UUID(g.Admin),
		Members:     UUIDs(g.Members),
		Name:        g.Name,
		Description: g.Description,
		GroupPhoto:  string(g.GroupPhoto),
		CreatedAt:   int64(g.CreatedAt),
	}
}
