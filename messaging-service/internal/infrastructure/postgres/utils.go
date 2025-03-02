package postgres

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func userIDs(arr pgtype.Array[uuid.UUID]) []domain.UserID {
	res := make([]domain.UserID, arr.Dims[0].Length)
	for i := range res {
		res[i] = domain.UserID(arr.Index(i).(uuid.UUID))
	}
	return res
}
