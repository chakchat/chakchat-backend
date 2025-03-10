package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/google/uuid"
)

var ErrNoCriteriaCpecified = errors.New("invalid input")

type GetUserRepository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUsersByCriteria(ctx context.Context, req storage.SearchUsersRequest) (*storage.SearchUsersResponse, error)
}

type GetRestrictionRepository interface {
	GetRestrictions(ctx context.Context, id uuid.UUID, field string) (*storage.FieldRestrictions, error)
}

type GetUserService struct {
	getUserRepo        GetUserRepository
	getRestrictionRepo GetRestrictionRepository
}

func NewGetService(getter GetUserRepository, restrictions GetRestrictionRepository) *GetUserService {
	return &GetUserService{
		getUserRepo:        getter,
		getRestrictionRepo: restrictions,
	}
}

func (g *GetUserService) GetUserByID(ctx context.Context, ownerId uuid.UUID, targetId uuid.UUID) (*models.User, error) {
	user, err := g.getUserRepo.GetUserById(ctx, targetId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if targetId == ownerId {
		return user, nil
	}

	if user.PhoneVisibility == models.RestrictionNone {
		user.Phone = ""
	}

	if user.PhoneVisibility == models.RestrictionSpecified {
		restr, err := g.getRestrictionRepo.GetRestrictions(ctx, targetId, "Phone")
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		if !canView(ownerId, restr.SpecifiedUsers) {
			user.Phone = ""
		}
	}

	if user.DateOfBirthVisibility == models.RestrictionNone {
		user.DateOfBirth = nil
	}

	if user.DateOfBirthVisibility == models.RestrictionSpecified {
		restr, err := g.getRestrictionRepo.GetRestrictions(ctx, targetId, "DateOfBirth")
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		if !canView(ownerId, restr.SpecifiedUsers) {
			user.DateOfBirth = nil
		}
	}

	return user, nil
}

func (g *GetUserService) GetUserByUsername(ctx context.Context, ownerId uuid.UUID, username string) (*models.User, error) {
	user, err := g.getUserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if user.ID == ownerId {
		return user, nil
	}

	if user.PhoneVisibility == models.RestrictionNone {
		user.Phone = ""
	}

	if user.PhoneVisibility == models.RestrictionSpecified {
		restr, err := g.getRestrictionRepo.GetRestrictions(ctx, user.ID, "Phone")
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		if !canView(ownerId, restr.SpecifiedUsers) {
			user.Phone = ""
		}
	}

	if user.DateOfBirthVisibility == models.RestrictionNone {
		user.DateOfBirth = nil
	}

	if user.DateOfBirthVisibility == models.RestrictionSpecified {
		restr, err := g.getRestrictionRepo.GetRestrictions(ctx, user.ID, "DateOfBirth")
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		if !canView(ownerId, restr.SpecifiedUsers) {
			user.DateOfBirth = nil
		}
	}

	return user, nil
}

func (g *GetUserService) GetUsersByCriteria(ctx context.Context, req storage.SearchUsersRequest) (*storage.SearchUsersResponse, error) {
	if req.Name == nil && req.Username == nil {
		return nil, ErrNoCriteriaCpecified
	}
	resp, err := g.getUserRepo.GetUsersByCriteria(ctx, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return resp, nil
}

func canView(ownerId uuid.UUID, specifiedUsers []uuid.UUID) bool {
	for _, id := range specifiedUsers {
		if id == ownerId {
			return true
		}
	}
	return false
}
