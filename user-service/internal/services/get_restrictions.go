package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/google/uuid"
)

type GetRestrictionsRepository interface {
	GetRestrictions(ctx context.Context, id uuid.UUID) (*models.UserRestrictions, error)
}

type GetRestrictionService struct {
	restrictionRepo GetRestrictionRepository
}

func NewGetRestrictionService(restrictionRepo GetRestrictionRepository) *GetRestrictionService {
	return &GetRestrictionService{
		restrictionRepo: restrictionRepo,
	}
}

func (g *GetRestrictionService) GetRestrictions(ctx context.Context, id uuid.UUID) (*models.UserRestrictions, error) {
	restr, err := g.restrictionRepo.GetRestriction(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return restr, nil
}
