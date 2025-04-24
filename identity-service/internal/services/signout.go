package services

import (
	"context"
	"fmt"

	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/google/uuid"
)

type RefreshTokenInvalidator interface {
	Invalidate(context.Context, jwt.Token) error
}

type SignOutService struct {
	invalidator   RefreshTokenInvalidator
	refreshConfig *jwt.Config
	deviceStorage DeviceStorage
}

func NewSignOutService(invalidator RefreshTokenInvalidator, refreshConf *jwt.Config, deviceStorage DeviceStorage) *SignOutService {
	return &SignOutService{
		invalidator:   invalidator,
		refreshConfig: refreshConf,
		deviceStorage: deviceStorage,
	}
}

func (s *SignOutService) SignOut(ctx context.Context, refresh jwt.Token) error {
	// idk should I check smth?
	if err := s.invalidator.Invalidate(ctx, refresh); err != nil {
		return fmt.Errorf("token invalidation failed: %s", err)
	}
	claims, err := jwt.Parse(s.refreshConfig, refresh)
	if err != nil {
		return fmt.Errorf("failed to parse refresh token: %s", err)
	}
	sub := claims[jwt.ClaimSub].(string)

	userID, err := uuid.Parse(sub)
	if err != nil {
		return fmt.Errorf("failed to parse sub claim")
	}

	if err := s.deviceStorage.Remove(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete device token: %s", err)
	}
	return nil
}
