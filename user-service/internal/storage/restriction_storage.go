package storage

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRestrictions struct {
	Phone       FieldRestriction
	DateOfBirth FieldRestriction
}

type FieldRestriction struct {
	OpenTo         string
	SpecifiedUsers []uuid.UUID
}

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
	return &restrictions, nil
}

func (s *RestrictionStorage) UpdateRestrictions(ctx context.Context, id uuid.UUID, restrictions UserRestrictions) (*models.UserRestrictions, error) {
	var userRestrictions models.UserRestrictions
	if err := s.db.WithContext(ctx).Where(&models.User{ID: id}).First(&userRestrictions).Error; err != nil {
		return nil, ErrNotFound
	}

	var phoneSpecifiedUsers []models.FieldRestrictionUser

	for _, id := range restrictions.Phone.SpecifiedUsers {
		phoneSpecifiedUsers = append(phoneSpecifiedUsers, models.FieldRestrictionUser{
			ID:                 uuid.New(),
			FieldRestrictionId: userRestrictions.Phone.ID,
			UserID:             id,
		})
	}

	var dateSpecifiedUsers []models.FieldRestrictionUser

	for _, id := range restrictions.DateOfBirth.SpecifiedUsers {
		dateSpecifiedUsers = append(dateSpecifiedUsers, models.FieldRestrictionUser{
			ID:                 uuid.New(),
			FieldRestrictionId: userRestrictions.DateOfBirth.ID,
			UserID:             id,
		})
	}
	if err := s.db.WithContext(ctx).Save(&models.UserRestrictions{UserId: id,
		Phone: models.FieldRestriction{
			OpenTo:         restrictions.Phone.OpenTo,
			SpecifiedUsers: phoneSpecifiedUsers,
		},
		DateOfBirth: models.FieldRestriction{
			OpenTo:         restrictions.DateOfBirth.OpenTo,
			SpecifiedUsers: dateSpecifiedUsers,
		}}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &userRestrictions, nil
}
