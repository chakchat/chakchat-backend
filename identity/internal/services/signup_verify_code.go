package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrSignUpKeyNotFound = errors.New("sign up key not found")
)

type SignUpMetaFindUpdater interface {
	FindMeta(ctx context.Context, signInKey uuid.UUID) (*SignUpMeta, bool, error)
	Store(context.Context, *SignUpMeta) error
}

type SignUpVerifyCodeService struct {
	storage SignUpMetaFindUpdater
}

func NewSignUpVerifyCodeService(storage SignUpMetaFindUpdater) *SignUpVerifyCodeService {
	return &SignUpVerifyCodeService{
		storage: storage,
	}
}

func (s *SignUpVerifyCodeService) VerifyCode(ctx context.Context, signUpKey uuid.UUID, code string) error {
	meta, ok, err := s.storage.FindMeta(ctx, signUpKey)
	if err != nil {
		return fmt.Errorf("sign up metadata finding failed: %s", err)
	}
	if !ok {
		return ErrSignUpKeyNotFound
	}
	if meta.Code != code {
		return ErrWrongCode
	}

	meta.Verified = true
	if err := s.storage.Store(ctx, meta); err != nil {
		return fmt.Errorf("sign up meta update failed: %s", err)
	}

	return nil
}
