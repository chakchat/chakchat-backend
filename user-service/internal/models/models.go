package models

import (
	"time"

	"github.com/google/uuid"
)

type Restriction string

const (
	RestrictionAll       = "everyone"
	RestrictionSpecified = "specified"
	RestrictionNone      = "only_me"
)

type User struct {
	ID          uuid.UUID
	Name        string
	Username    string
	Phone       string
	DateOfBirth *time.Time
	PhotoURL    string
	CreatedAt   int64

	DateOfBirthVisibility Restriction `gorm:"default:everyone"`
	PhoneVisibility       Restriction `gorm:"default:everyone"`
}

type FieldRestriction struct {
	OwnerID        uuid.UUID `gorm:"primaryKey"`
	FieldName      string
	SpecifiedUsers []uuid.UUID
}

type FieldRestrictionUser struct {
	UserID             uuid.UUID         `gorm:"type:uuid"`
	FieldRestrictionId uuid.UUID         `gorm:"type:uuid"`
	FieldRestriction   *FieldRestriction `gorm:"constraint:OnDelete:CASCADE"`
}
