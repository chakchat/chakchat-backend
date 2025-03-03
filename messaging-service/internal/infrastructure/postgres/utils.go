package postgres

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

func userIDs(arr []uuid.UUID) []domain.UserID {
	res := make([]domain.UserID, len(arr))
	for i := range res {
		res[i] = domain.UserID(arr[i])
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
