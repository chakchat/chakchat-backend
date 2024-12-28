package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/chakchat/chakchat/backend/identity/internal/userservice"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrSendCodeFreqExceeded = errors.New("send code operation frequency exceeded")
)

type SmsSender interface {
	SendSms(phone string, message string) error
}

type SignInMeta struct {
	SignInID    uuid.UUID
	LastRequest time.Time
	Code        string

	UserId   uuid.UUID
	Username string
}

type SignInMetaStorer interface {
	Store(SignInMeta) error
}

type CodeSender struct {
	sms    SmsSender
	users  userservice.UsersServiceClient
	storer SignInMetaStorer
}

func NewCodeSender(sms SmsSender, storer SignInMetaStorer, users userservice.UsersServiceClient) *CodeSender {
	return &CodeSender{
		sms:    sms,
		users:  users,
		storer: storer,
	}
}

func (s *CodeSender) SendCode(ctx context.Context, phone string) (signInKey uuid.UUID, err error) {
	meta := SignInMeta{
		SignInID:    uuid.New(),
		LastRequest: time.Now().UTC(),
	}

	user, err := s.users.GetUser(ctx, &userservice.UserRequest{
		PhoneNumber: phone,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("user gRPC call error: %s", err)
	}
	if user.Status != userservice.UserResponseStatus_SUCCESS {
		switch user.Status {
		case userservice.UserResponseStatus_FAILED:
			return uuid.Nil, errors.New("unknown gRPC GetUser() error")
		case userservice.UserResponseStatus_NOT_FOUND:
			return uuid.Nil, ErrUserNotFound
		default:
			return uuid.Nil, errors.New("implement one another status handling")
		}
	}

	meta.UserId = uuid.MustParse(user.UserId.GetValue())
	meta.Username = *user.UserName
	meta.Code = genCode()

	if err := s.sms.SendSms(phone, meta.Code); err != nil {
		return uuid.Nil, fmt.Errorf("send SMS error: %s", err)
	}

	if err := s.storer.Store(meta); err != nil {
		return uuid.Nil, fmt.Errorf("storage error: %s", err)
	}

	return meta.SignInID, err
}

func genCode() string {
	n := 100000 + rand.IntN(900000)
	return strconv.Itoa(n)
}
