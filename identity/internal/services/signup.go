package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat/backend/identity/internal/jwt"
	"github.com/chakchat/chakchat/backend/identity/internal/userservice"
	"github.com/google/uuid"
)

var (
	ErrUsernameAlreadyExists     = errors.New("username already exists")
	ErrCreateUserValidationError = errors.New("create user validation error")
	ErrPhoneNotVerified          = errors.New("phone is not verified")
)

type CreateUserData struct {
	Username string
	Name     string
}

type SignUpMetaFindRemover interface {
	FindMeta(ctx context.Context, signInKey uuid.UUID) (*SignUpMeta, bool, error)
	Remove(ctx context.Context, signInKey uuid.UUID) error
}

type SignUpService struct {
	accessConf  *jwt.Config
	refreshConf *jwt.Config

	users   userservice.UserServiceClient
	storage SignUpMetaFindRemover
}

func NewSignUpService(users userservice.UserServiceClient, storage SignUpMetaFindRemover) *SignUpService {
	return &SignUpService{
		users:   users,
		storage: storage,
	}
}

func (s *SignUpService) SignUp(ctx context.Context, signUpKey uuid.UUID, user CreateUserData) (jwt.Pair, error) {
	meta, err := s.checkMeta(ctx, signUpKey)
	if err != nil {
		return jwt.Pair{}, err
	}

	userResp, err := s.createUser(ctx, meta.Phone, user.Name, user.Username)
	if err != nil {
		return jwt.Pair{}, err
	}

	tokens, err := s.generateTokens(userResp.UserId, *userResp.Name, *userResp.UserName)
	if err != nil {
		return jwt.Pair{}, err
	}

	if err := s.storage.Remove(ctx, signUpKey); err != nil {
		return jwt.Pair{}, fmt.Errorf("sign up meta removal failed: %s", err)
	}

	return tokens, nil
}

func (s *SignUpService) generateTokens(id *userservice.UUID, name, username string) (jwt.Pair, error) {
	claims := jwt.Claims{
		jwt.ClaimSub:      id,
		jwt.ClaimName:     name,
		jwt.ClaimUsername: username,
	}

	var pair jwt.Pair
	var err error
	if pair.Access, err = jwt.Generate(s.accessConf, claims); err != nil {
		return jwt.Pair{}, fmt.Errorf("access token generation failed: %s", err)
	}
	if pair.Refresh, err = jwt.Generate(s.refreshConf, claims); err != nil {
		return jwt.Pair{}, fmt.Errorf("refresh token generation failed: %s", err)
	}

	return pair, nil
}

func (s *SignUpService) checkMeta(ctx context.Context, signUpKey uuid.UUID) (*SignUpMeta, error) {
	meta, ok, err := s.storage.FindMeta(ctx, signUpKey)
	if err != nil {
		return nil, fmt.Errorf("sign up meta finding failed: %s", err)
	}

	if !ok {
		return nil, ErrSignUpKeyNotFound
	}
	if !meta.Verified {
		return nil, ErrPhoneNotVerified
	}

	return meta, nil
}

func (s *SignUpService) createUser(ctx context.Context, phone, name, username string) (*userservice.CreateUserResponse, error) {
	resp, err := s.users.CreateUser(ctx, &userservice.CreateUserRequest{
		PhoneNumber: phone,
		Name:        name,
		Username:    username,
	})
	if err != nil {
		return nil, fmt.Errorf("create user gRPC call failed: %s", err)
	}

	switch resp.Status {
	case userservice.CreateUserStatus_CREATED:
		return resp, nil
	case userservice.CreateUserStatus_ALREADY_EXISTS:
		// Actually, I'm not sure that only username may cause this error
		return nil, ErrUsernameAlreadyExists
	case userservice.CreateUserStatus_CREATE_FAILED:
		return nil, errors.New("user creation failed")
	default:
		return nil, errors.New("unknown response status")
	}
}
