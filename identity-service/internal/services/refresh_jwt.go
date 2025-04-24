package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/google/uuid"
)

var (
	ErrRefreshTokenExpired     = errors.New("refresh token expired")
	ErrRefreshTokenInvalidated = errors.New("refresh token invalidated")
	ErrInvalidJWT              = errors.New("jwt token is invalid")
	ErrInvalidTokenType        = errors.New("jwt token is invalid")
)

type RefreshTokenCheckInvalidator interface {
	Invalidated(context.Context, jwt.Token) (bool, error)
	Invalidate(context.Context, jwt.Token) error
}

type RefreshService struct {
	accessConf    *jwt.Config
	refreshConf   *jwt.Config
	deviceStorage DeviceStorage
	checker       RefreshTokenCheckInvalidator
}

func NewRefreshService(checker RefreshTokenCheckInvalidator, accessConf, refreshConf *jwt.Config, deviceStorage DeviceStorage) *RefreshService {
	return &RefreshService{
		accessConf:    accessConf,
		deviceStorage: deviceStorage,
		refreshConf:   refreshConf,
		checker:       checker,
	}
}

func (s *RefreshService) Refresh(ctx context.Context, refresh jwt.Token) (jwt.Pair, error) {
	if err := s.validate(ctx, refresh); err != nil {
		return jwt.Pair{}, err
	}

	parsed, err := jwt.Parse(s.refreshConf, refresh)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			return jwt.Pair{}, ErrRefreshTokenExpired
		}
		if err == jwt.ErrInvalidTokenType {
			return jwt.Pair{}, ErrInvalidTokenType
		}
		return jwt.Pair{}, ErrInvalidJWT
	}

	if err := s.checker.Invalidate(ctx, refresh); err != nil {
		return jwt.Pair{}, fmt.Errorf("refresh token invalidation failed: %s", err)
	}

	claims := extractPublic(parsed)

	sub := claims[jwt.ClaimSub].(string)

	userId, err := uuid.Parse(sub)
	if err != nil {
		return jwt.Pair{}, fmt.Errorf("failed to parse sub claim")
	}
	s.deviceStorage.Refresh(ctx, userId)

	var pair jwt.Pair
	if pair.Access, err = jwt.Generate(s.accessConf, claims); err != nil {
		return jwt.Pair{}, fmt.Errorf("access token generation failed: %s", err)
	}
	if pair.Refresh, err = jwt.Generate(s.refreshConf, claims); err != nil {
		return jwt.Pair{}, fmt.Errorf("refresh token generation failed: %s", err)
	}

	return pair, nil
}

func extractPublic(claims jwt.Claims) jwt.Claims {
	return jwt.Claims{
		jwt.ClaimSub:      claims[jwt.ClaimSub],
		jwt.ClaimName:     claims[jwt.ClaimName],
		jwt.ClaimUsername: claims[jwt.ClaimUsername],
	}
}

func (s *RefreshService) validate(ctx context.Context, refresh jwt.Token) error {
	invalidated, err := s.checker.Invalidated(ctx, refresh)
	if err != nil {
		return err
	}
	if invalidated {
		return ErrRefreshTokenInvalidated
	}
	return nil
}
