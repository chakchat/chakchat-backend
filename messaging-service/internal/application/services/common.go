package services

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

func GetSecondUserSlice(users [2]domain.UserID, first domain.UserID) []uuid.UUID {
	var second domain.UserID
	if users[0] == first {
		second = users[1]
	} else {
		second = users[0]
	}
	return []uuid.UUID{uuid.UUID(second)}
}
