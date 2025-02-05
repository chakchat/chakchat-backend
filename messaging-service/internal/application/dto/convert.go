package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

func UserIDs(users []uuid.UUID) []domain.UserID {
	res := make([]domain.UserID, len(users))
	for i, u := range users {
		res[i] = domain.UserID(u)
	}
	return res
}

func UUIDs(users []domain.UserID) []uuid.UUID {
	res := make([]uuid.UUID, len(users))
	for i, u := range users {
		res[i] = uuid.UUID(u)
	}
	return res
}
