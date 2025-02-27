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
	GetRestriction(ctx context.Context, id uuid.UUID) (*models.UserRestrictions, error)
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
	user, err := g.getUserRepo.GetUserById(ctx, ownerId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	restr, err := g.getRestrictionRepo.GetRestriction(ctx, ownerId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if !canView(ownerId, restr.Phone) {
		user.Phone = ""
	}

	if !canView(ownerId, restr.DateOfBirth) {
		user.DateOfBirth = nil
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

	restr, err := g.getRestrictionRepo.GetRestriction(ctx, user.ID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if !canView(ownerId, restr.Phone) {
		user.Phone = ""
	}

	if !canView(ownerId, restr.DateOfBirth) {
		user.DateOfBirth = nil
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

func canView(ownerId uuid.UUID, restr models.FieldRestriction) bool {
	switch restr.OpenTo {
	case "everyone":
		return true
	case "only_me":
		return ownerId == restr.SpecifiedUsers[0].ID
	case "specified":
		for _, id := range restr.SpecifiedUsers {
			if id.ID == ownerId {
				return true
			}
		}
		return false
	default:
		return false
	}
}
