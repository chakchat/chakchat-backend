package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Username    string
	Name        string
	Phone       string
	DateOfBirth *time.Time
	PhotoURL    string //'' if no photo
	CreatedAt   int64
}

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) GetUser(ctx context.Context, phone string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).Where(&User{Phone: phone}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &User{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		PhotoURL:    user.PhotoURL,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *UserStorage) CreateUser(ctx context.Context, user *User) (*User, error) {
	var newUser User = User{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		PhotoURL:    user.PhotoURL,
		CreatedAt:   user.CreatedAt,
	}

	if err := s.db.WithContext(ctx).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		return nil, gorm.ErrDuplicatedKey
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, err
	}
	return &newUser, nil
}
