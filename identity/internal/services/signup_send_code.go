package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chakchat/chakchat/backend/identity/internal/userservice"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type SignUpMeta struct {
	SignUpKey   uuid.UUID
	LastRequest time.Time
	Phone       string
	Code        string
	Verified    bool
}

type SignUpMetaFindStorer interface {
	FindMetaByPhone(ctx context.Context, phone string) (*SignUpMeta, bool, error)
	Store(context.Context, *SignUpMeta) error
}

type SignUpSendCodeService struct {
	config *CodeConfig

	sms     SmsSender
	storage SignUpMetaFindStorer
	users   userservice.UserServiceClient
}

func NewSignUpSendCodeService(config *CodeConfig, sms SmsSender,
	storage SignUpMetaFindStorer, users userservice.UserServiceClient) *SignUpSendCodeService {
	return &SignUpSendCodeService{
		config:  config,
		sms:     sms,
		storage: storage,
		users:   users,
	}
}

func (s *SignUpSendCodeService) SendCode(ctx context.Context, phone string) (signUpKey uuid.UUID, err error) {
	if err := s.validateSendFreq(ctx, phone); err != nil {
		return uuid.Nil, err
	}

	if err := s.validateUserExistence(ctx, phone); err != nil {
		return uuid.Nil, err
	}

	meta := SignUpMeta{
		SignUpKey:   uuid.New(),
		LastRequest: nowUTC(),
		Phone:       phone,
		Code:        genCode(),
		Verified:    false,
	}

	sms := renderCodeSMS(meta.Code)
	if err := s.sms.SendSms(ctx, phone, sms); err != nil {
		return uuid.Nil, fmt.Errorf("send SMS error: %s", err)
	}

	if err := s.storage.Store(ctx, &meta); err != nil {
		return uuid.Nil, fmt.Errorf("storage error: %s", err)
	}

	return meta.SignUpKey, err
}

func (s *SignUpSendCodeService) validateUserExistence(ctx context.Context, phone string) error {
	// While testing I found out that first request takes exactly 8 seconds more
	user, err := s.users.GetUser(ctx, &userservice.UserRequest{
		PhoneNumber: phone,
	})

	if err != nil {
		return fmt.Errorf("user gRPC call error: %s", err)
	}

	if user.Status == userservice.UserResponseStatus_NOT_FOUND {
		return nil
	}

	if user.Status == userservice.UserResponseStatus_SUCCESS {
		return ErrUserAlreadyExists
	}
	if user.Status == userservice.UserResponseStatus_FAILED {
		return errors.New("unknown gRPC GetUser() error")
	}
	return fmt.Errorf("unexpected user service status: %v", user.Status)
}

func (s *SignUpSendCodeService) validateSendFreq(ctx context.Context, phone string) error {
	prevMeta, ok, err := s.storage.FindMetaByPhone(ctx, phone)
	if err != nil {
		return fmt.Errorf("finding SignUpMeta error: %s", err)
	}
	if ok && prevMeta.LastRequest.Add(s.config.SendFrequency).Compare(nowUTC()) > 0 {
		return ErrSendCodeFreqExceeded
	}
	return nil
}
