package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

type UserRepository interface {
	// Returns NotFound error if not found.
	GetUserByPhone(ctx context.Context, phone string) (*storage.User, error)
	CreateUser(ctx context.Context, user *storage.User) (*storage.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewGetFileService(userHandler UserRepository) *UserService {
	return &UserService{
		userRepo: userHandler,
	}
}

func (s *UserService) GetUser(ctx context.Context, phone string) (*storage.User, error) {
	user, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &storage.User{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		PhotoURL:    user.PhotoURL,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *storage.User) (*storage.User, error) {
	newUser, err := s.userRepo.CreateUser(ctx, user)
	if errors.Is(err, storage.ErrAlreadyExists) {
		return nil, ErrAlreadyExists
	}
	return newUser, nil
}
