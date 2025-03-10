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
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username    string
	Name        string
	Phone       string
	DateOfBirth *time.Time
	PhotoURL    string
	CreatedAt   int64 `gorm:"autoCreateTime"`

	DateOfBirthVisibility Restriction `gorm:"default:everyone"`
	PhoneVisibility       Restriction `gorm:"default:everyone"`
}

type FieldRestriction struct {
	OwnerID        uuid.UUID `gorm:"primaryKey"`
	FieldName      string
	SpecifiedUsers []FieldRestrictionUser `gorm:"foreignKey:FieldRestrictionId;constraint:OnDelete:Cascade"`
}

type FieldRestrictionUser struct {
	UserID             uuid.UUID         `gorm:"type:uuid"`
	FieldRestrictionId uuid.UUID         `gorm:"type:uuid"`
	FieldRestriction   *FieldRestriction `gorm:"constraint:OnDelete:CASCADE"`
}
