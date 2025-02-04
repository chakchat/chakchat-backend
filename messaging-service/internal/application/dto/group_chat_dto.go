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
	members := make([]uuid.UUID, len(g.Members))
	for i, u := range g.Members {
		members[i] = uuid.UUID(u)
	}

	return GroupChatDTO{
		ID:          uuid.UUID(g.ChatID),
		Admin:       uuid.UUID(g.Admin),
		Members:     members,
		Name:        g.Name,
		Description: g.Description,
		GroupPhoto:  string(g.GroupPhoto),
		CreatedAt:   int64(g.CreatedAt),
	}
}
