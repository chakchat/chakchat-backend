package services

import (
	"context"
	"fmt"

	"github.com/chakchat/chakchat/backend/identity/internal/jwt"
)

type RefreshTokenInvalidator interface {
	Invalidate(context.Context, jwt.Token) error
}

type SignOutService struct {
	invalidator RefreshTokenInvalidator
}

func NewSignOutService(invalidator RefreshTokenInvalidator) *SignOutService {
	return &SignOutService{
		invalidator: invalidator,
	}
}

func (s *SignOutService) SignOut(ctx context.Context, refresh jwt.Token) error {
	// idk should I check smth?
	if err := s.invalidator.Invalidate(ctx, refresh); err != nil {
		return fmt.Errorf("token invalidation failed: %s", err)
	}
	return nil
}
