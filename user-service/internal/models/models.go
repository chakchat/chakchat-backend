package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username    string
	Name        string
	Phone       string
	DateOfBirth *time.Time
	PhotoURL    string
	CreatedAt   int64
}

type UserRestrictions struct {
	UserId      uuid.UUID        `gorm:"primaryKey"`
	Phone       FieldRestriction `gorm:"embedded"`
	DateOfBirth FieldRestriction `gorm:"embedded"`
}

type FieldRestriction struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	OpenTo         string
	SpecifiedUsers []FieldRestrictionUser `gorm:"foreignKey:UserID"`
}

type FieldRestrictionUser struct {
	ID                 uuid.UUID `gorm:"primaryKey"`
	FieldRestrictionId uuid.UUID
	UserID             uuid.UUID
}
