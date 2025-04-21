package services

import (
	"context"

	"github.com/google/uuid"
)

type GetRestrictionsRepository interface {
	GetAllowedUserIDs(ctx context.Context, id uuid.UUID, field string) ([]uuid.UUID, error)
}

type GetRestrictionService struct {
	restrictionRepo GetRestrictionsRepository
}

func NewGetRestrictionService(restrictionRepo GetRestrictionsRepository) *GetRestrictionService {
	return &GetRestrictionService{
		restrictionRepo: restrictionRepo,
	}
}

func (g *GetRestrictionService) GetAllowedUserIDs(ctx context.Context, id uuid.UUID, field string) ([]uuid.UUID, error) {
	restr, err := g.restrictionRepo.GetAllowedUserIDs(ctx, id, field)
	if err != nil {
		return nil, err
	}

	return restr, nil
}
