package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/google/uuid"
)

type GetRestrictionsRepository interface {
	GetRestrictions(ctx context.Context, id uuid.UUID, field string) (*storage.FieldRestrictions, error)
}

type GetRestrictionService struct {
	restrictionRepo GetRestrictionsRepository
}

func NewGetRestrictionService(restrictionRepo GetRestrictionsRepository) *GetRestrictionService {
	return &GetRestrictionService{
		restrictionRepo: restrictionRepo,
	}
}

func (g *GetRestrictionService) GetRestrictions(ctx context.Context, id uuid.UUID, field string) (*storage.FieldRestrictions, error) {
	restr, err := g.restrictionRepo.GetRestrictions(ctx, id, field)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return restr, nil
}
