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

type SignUpMetaFindRemover interface {
	FindMeta(ctx context.Context, signInKey uuid.UUID) (*SignUpMeta, bool, error)
	Remove(ctx context.Context, signInKey uuid.UUID) error
}

type SignUpVerifyCodeService struct {
	storage SignUpMetaFindRemover
}

func NewSignUpVerifyCodeService(storage SignUpMetaFindRemover) *SignUpVerifyCodeService {
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

	if err := s.storage.Remove(ctx, signUpKey); err != nil {
		return fmt.Errorf("sign up key removal failed: %s", err)
	}

	return nil
}
