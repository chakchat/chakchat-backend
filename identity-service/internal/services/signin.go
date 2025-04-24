package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/google/uuid"
)

var (
	ErrSignInKeyNotFound = errors.New("sign in key not found")
	ErrWrongCode         = errors.New("wrong phone verification code")
)

type DeviceInfo struct {
	Type        string
	DeviceToken string
}

type SignInMetaFindRemover interface {
	FindMeta(ctx context.Context, signInKey uuid.UUID) (*SignInMeta, bool, error)
	Remove(ctx context.Context, signInKey uuid.UUID) error
}

type DeviceStorage interface {
	Store(ctx context.Context, userID uuid.UUID, info *DeviceInfo) error
	Refresh(ctx context.Context, userID uuid.UUID) error
	Remove(ctx context.Context, userID uuid.UUID) error
}
type SignInService struct {
	storage       SignInMetaFindRemover
	deviceStorage DeviceStorage
	accessConf    *jwt.Config
	refreshConf   *jwt.Config
}

func NewSignInService(storage SignInMetaFindRemover, accessConf, refreshConf *jwt.Config, deviceStorage DeviceStorage) *SignInService {
	return &SignInService{
		storage:       storage,
		deviceStorage: deviceStorage,
		accessConf:    accessConf,
		refreshConf:   refreshConf,
	}
}

func (s *SignInService) SignIn(ctx context.Context, signInKey uuid.UUID, code string, device *DeviceInfo) (jwt.Pair, error) {
	meta, ok, err := s.storage.FindMeta(ctx, signInKey)
	if err != nil {
		return jwt.Pair{}, fmt.Errorf("sign in metadata finding failed: %s", err)
	}
	if !ok {
		return jwt.Pair{}, ErrSignInKeyNotFound
	}
	if meta.Code != code {
		return jwt.Pair{}, ErrWrongCode
	}

	claims := jwt.Claims{
		jwt.ClaimSub:      meta.UserId,
		jwt.ClaimName:     meta.Name,
		jwt.ClaimUsername: meta.Username,
	}
	var pair jwt.Pair
	if pair.Access, err = jwt.Generate(s.accessConf, claims); err != nil {
		return jwt.Pair{}, fmt.Errorf("access token generation failed: %s", err)
	}
	if pair.Refresh, err = jwt.Generate(s.refreshConf, claims); err != nil {
		return jwt.Pair{}, fmt.Errorf("refresh token generation failed: %s", err)
	}

	if err := s.storage.Remove(ctx, signInKey); err != nil {
		return jwt.Pair{}, fmt.Errorf("sign in key removal failed: %s", err)
	}

	if device != nil {
		if err := s.deviceStorage.Store(ctx, meta.UserId, device); err != nil {
			return jwt.Pair{}, fmt.Errorf("failed to store device info: %s", err)
		}
	}

	return pair, nil
}
