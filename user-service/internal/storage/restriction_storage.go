package storage

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestrictionStorage struct {
	db *gorm.DB
}

func NewRestrictionStorage(db *gorm.DB) *RestrictionStorage {
	return &RestrictionStorage{
		db: db,
	}
}

func (s *RestrictionStorage) GetRestriction(ctx context.Context, id uuid.UUID) (*models.UserRestrictions, error) {
	var restrictions models.UserRestrictions
	if err := s.db.WithContext(ctx).First(&restrictions, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &models.UserRestrictions{
		UserId:      restrictions.UserId,
		Phone:       restrictions.Phone,
		DateOfBirth: restrictions.DateOfBirth,
	}, nil
}
