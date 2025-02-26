package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRestrictions struct {
	UserId      uuid.UUID        `gorm:"primaryKey"`
	Phone       FieldRestriction `gorm:"embedded"`
	DateOfBirth FieldRestriction `gorm:"embedded"`
}

type FieldRestriction struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	OpenTo         string
	SpecifiedUsers []User `gorm:"foreignKey:ID"`
}

type RestrictionStorage struct {
	db *gorm.DB
}

func NewRestrictionStorage(db *gorm.DB) *RestrictionStorage {
	return &RestrictionStorage{
		db: db,
	}
}

func (s *RestrictionStorage) GetRestriction(ctx context.Context, id uuid.UUID) (*UserRestrictions, error) {
	var restrictions UserRestrictions
	if err := s.db.WithContext(ctx).First(&restrictions, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &UserRestrictions{
		UserId:      restrictions.UserId,
		Phone:       restrictions.Phone,
		DateOfBirth: restrictions.DateOfBirth,
	}, nil
}
