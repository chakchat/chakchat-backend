package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/google/uuid"
)

type UpdateRestrictionsRepo interface {
	UpdateRestrictions(ctx context.Context, id uuid.UUID, phone storage.FieldRestriction, date storage.FieldRestriction) (*models.UserRestrictions, error)
}

type UpdateRestrictionsService struct {
	repo UpdateRestrictionsRepo
}

func NewUpdateRestrService(repo UpdateRestrictionsRepo) *UpdateRestrictionsService {
	return &UpdateRestrictionsService{
		repo: repo,
	}
}

func (s *UpdateRestrictionsService) UpdateRestrictions(ctx context.Context, id uuid.UUID, phone storage.FieldRestriction, date storage.FieldRestriction) (*models.UserRestrictions, error) {
	updatedRestr, err := s.repo.UpdateRestrictions(ctx, id, phone, date)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrValidationError
		}
		return nil, err
	}
	return updatedRestr, nil
}
