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

func uuids(s []domain.UserID) []uuid.UUID {
	res := make([]uuid.UUID, len(s))
	for i := range res {
		res[i] = uuid.UUID(s[i])
	}
	return res
}

func sliceMisses[T comparable](orig, comp []T) []T {
	compMap := make(map[T]bool, len(comp))
	for _, t := range comp {
		compMap[t] = true
	}

	var misses []T

	for _, t := range orig {
		if !compMap[t] {
			misses = append(misses, t)
		}
	}

	return misses
}
