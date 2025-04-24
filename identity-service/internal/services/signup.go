package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/identity-service/internal/userservice"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
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

	users         userservice.UserServiceClient
	storage       SignUpMetaFindRemover
	deviceStorage DeviceStorage
}

func NewSignUpService(accessConf *jwt.Config, refreshConf *jwt.Config, users userservice.UserServiceClient,
	storage SignUpMetaFindRemover) *SignUpService {
	return &SignUpService{
		accessConf:  accessConf,
		refreshConf: refreshConf,
		users:       users,
		storage:     storage,
	}
}

func (s *SignUpService) SignUp(ctx context.Context, signUpKey uuid.UUID, user CreateUserData, device *DeviceInfo) (jwt.Pair, error) {
	meta, err := s.checkMeta(ctx, signUpKey)
	if err != nil {
		return jwt.Pair{}, err
	}

	userResp, err := s.createUser(ctx, meta.Phone, user.Name, user.Username)
	if err != nil {
		return jwt.Pair{}, err
	}

	id, err := uuid.Parse(userResp.UserId.Value)
	if err != nil {
		return jwt.Pair{}, fmt.Errorf("failed to parse UserID from user service: %s", err)
	}

	claims := jwt.Claims{
		jwt.ClaimSub:      id,
		jwt.ClaimName:     userResp.Name,
		jwt.ClaimUsername: userResp.UserName,
	}

	var tokens jwt.Pair
	if tokens.Access, err = jwt.Generate(s.accessConf, claims); err != nil {
		return jwt.Pair{}, fmt.Errorf("access token generation failed: %s", err)
	}
	if tokens.Refresh, err = jwt.Generate(s.refreshConf, claims); err != nil {
		return jwt.Pair{}, fmt.Errorf("refresh token generation failed: %s", err)
	}

	if err := s.storage.Remove(ctx, signUpKey); err != nil {
		return jwt.Pair{}, fmt.Errorf("sign up meta removal failed: %s", err)
	}

	userId, err := uuid.Parse(userResp.UserId.Value)
	if err != nil {
		return jwt.Pair{}, fmt.Errorf("can't parse user id")
	}
	if device != nil {
		if err := s.deviceStorage.Store(ctx, userId, device); err != nil {
			return jwt.Pair{}, fmt.Errorf("failed to store device info: %s", err)
		}
	}

	return tokens, nil
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
