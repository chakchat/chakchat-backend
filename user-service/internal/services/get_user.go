package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/google/uuid"
)

var ErrNoCriteriaCpecified = errors.New("invalid input")

type GetUserRepository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*storage.User, error)
	GetUserByUsername(ctx context.Context, username string) (*storage.User, error)
	GetUsersByCriteria(ctx context.Context, req storage.SearchUsersRequest) (*storage.SearchUsersResponse, error)
}

type GetRestrictionRepository interface {
	GetRestriction(ctx context.Context, id uuid.UUID) (*storage.UserRestrictions, error)
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

func (g *GetUserService) GetUserByID(ctx context.Context, ownerId uuid.UUID, targetId uuid.UUID) (*storage.User, error) {
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
		user.Phone = nil
	}

	if !canView(ownerId, restr.DateOfBirth) {
		user.DateOfBirth = nil
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

func (g *GetUserService) GetUserByUsername(ctx context.Context, ownerId uuid.UUID, username string) (*storage.User, error) {
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
		user.Phone = nil
	}

	if !canView(ownerId, restr.DateOfBirth) {
		user.DateOfBirth = nil
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
	return &storage.SearchUsersResponse{
		Users: resp.Users,
		Page:  resp.Page,
		Count: resp.Count,
	}, nil
}

func canView(ownerId uuid.UUID, restr storage.FieldRestriction) bool {
	switch restr.OpenTo {
	case "everyone":
		return true
	case "only_me":
		return ownerId == restr.SpecifiedUsers[0]
	case "specified":
		for _, id := range restr.SpecifiedUsers {
			if id == ownerId {
				return true
			}
		}
		return false
	default:
		return false
	}
}
