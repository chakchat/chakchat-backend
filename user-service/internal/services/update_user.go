package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
)

var ErrValidationError = errors.New("invalid input")

type UpdateUserRepository interface {
	UpdateUser(ctx context.Context, user *models.User, req *storage.UpdateUserRequest) (*models.User, error)
}

type UpdateUserService struct {
	updateRepo UpdateUserRepository
}

func NewUpdateUserService(updateRepo UpdateUserRepository) *UpdateUserService {
	return &UpdateUserService{
		updateRepo: updateRepo,
	}
}

func (u *UpdateUserService) UpdateUser(ctx context.Context, user *models.User, req *storage.UpdateUserRequest) (*models.User, error) {
	updatedUser, err := u.updateRepo.UpdateUser(ctx, user, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrValidationError
		}
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, ErrValidationError
		}
		return nil, err
	}
	return updatedUser, nil
}
