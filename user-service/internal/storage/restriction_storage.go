package storage

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldRestrictions struct {
	Field          string
	OpenTo         models.Restriction
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

func (s *RestrictionStorage) GetRestrictions(ctx context.Context, id uuid.UUID, field string) (*FieldRestrictions, error) {
	var fieldRestriction models.FieldRestriction
	var fieldSpecifiedUsers []uuid.UUID

	if err := s.db.WithContext(ctx).Where("owner_id = ? AND field_name = ?", id, field).Preload("SpecifiedUsers").Find(&fieldRestriction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	for _, user := range fieldRestriction.SpecifiedUsers {
		fieldSpecifiedUsers = append(fieldSpecifiedUsers, user.UserID)
	}

	return &FieldRestrictions{
		Field:          field,
		SpecifiedUsers: fieldSpecifiedUsers,
	}, nil
}

func (s *RestrictionStorage) UpdateRestrictions(ctx context.Context, id uuid.UUID, restrictions FieldRestrictions) (*FieldRestrictions, error) {

	var user models.User

	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if restrictions.Field == "Phone" {
		user.PhoneVisibility = restrictions.OpenTo
	} else {
		user.DateOfBirthVisibility = restrictions.OpenTo
	}

	var add []uuid.UUID
	var del []uuid.UUID
	_, err := s.GetRestrictions(ctx, id, restrictions.Field)
	if err != nil && err != gorm.ErrUnsupportedRelation {
		return nil, err
	}
	if err == nil {
		phoneRestr, err := s.GetRestrictions(ctx, id, restrictions.Field)
		if err != nil {
			return nil, err
		}

		add = recordMisses(phoneRestr.SpecifiedUsers, restrictions.SpecifiedUsers)
		del = recordMisses(restrictions.SpecifiedUsers, phoneRestr.SpecifiedUsers)

		err = s.db.WithContext(ctx).Where("field_restriction_id = ?", id).Delete(&models.FieldRestrictionUser{}, del).Error
		if err != nil {
			return nil, err
		}
	} else {
		fieldRestriction := models.FieldRestriction{
			OwnerID:   id,
			FieldName: restrictions.Field,
		}
		if err := s.db.Create(&fieldRestriction).Error; err != nil {
			return nil, err
		}
		add = restrictions.SpecifiedUsers
	}

	for _, record := range add {
		specifiedUser := models.FieldRestrictionUser{
			UserID:             record,
			FieldRestrictionId: id,
		}
		if err := s.db.Create(&specifiedUser).Error; err != nil {
			return nil, err
		}
	}

	return &restrictions, nil
}

func recordMisses(orig, comp []uuid.UUID) []uuid.UUID {
	compMap := make(map[uuid.UUID]bool, len(comp))
	for _, t := range comp {
		compMap[t] = true
	}

	var misses []uuid.UUID

	for _, t := range orig {
		if !compMap[t] {
			misses = append(misses, t)
		}
	}

	return misses
}
