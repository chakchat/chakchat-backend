package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("repository: not found")
var ErrAlreadyExists = errors.New("already exists")

type User struct {
	ID          uuid.UUID
	Username    string
	Name        string
	Phone       string
	DateOfBirth *time.Time
	PhotoURL    string //'' if no photo
	CreatedAt   int64
}

type UserRepository interface {
	// Returns NotFound error if not found.
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, bool, error)
}

type UserService struct {
	user_handler UserRepository
}

func NewGetFileService(user_handler UserRepository) *UserService {
	return &UserService{
		user_handler: user_handler,
	}
}

func (s *UserService) GetUser(ctx context.Context, phone string) (*User, error) {
	user, err := s.user_handler.GetUserByPhone(ctx, phone)
	if user == nil && err == nil {
		return nil, ErrNotFound
	} else if user == nil {
		return nil, err
	} else {
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
}

func (s *UserService) CreateUser(ctx context.Context, user *User) (*User, error) {
	newUser, duplicated, err := s.user_handler.CreateUser(ctx, user)
	if duplicated {
		return nil, ErrAlreadyExists
	} else if newUser == nil {
		return nil, err
	}
	return newUser, nil

}
