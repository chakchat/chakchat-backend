package services

import (
	"context"
	"testing"
	"time"

	"github.com/chakchat/chakchat-backend/identity-service/internal/userservice"
	"github.com/google/uuid"
	smsaero_golang "github.com/smsaero/smsaero_golang/smsaero"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_Success(t *testing.T) {
	// Arrange
	userService := userServiceMock{
		resp: &userservice.UserResponse{
			Status:   userservice.UserResponseStatus_SUCCESS,
			Name:     new(string),
			UserName: new(string),
			UserId: &userservice.UUID{
				Value: "6c056fb3-7efc-483a-9506-4336456ac79f",
			},
		},
	}
	metaStorage := metaStorageFake{}
	smsSender := smsStub{}
	config := CodeConfig{
		SendFrequency: 1 * time.Minute,
	}

	sender := NewSignInSendCodeService(&config, smsSender, &metaStorage, userService)

	// Act
	signInKey, err := sender.SendCode(context.Background(), "+79998887766")

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, signInKey)
}

func Test_Error(t *testing.T) {
	// Arrange
	now, _ := time.Parse(time.DateTime, "2024-12-28 03:00:01")
	nowUTC = func() time.Time { return now }

	userService := &userServiceMock{
		resp: &userservice.UserResponse{
			Status:   userservice.UserResponseStatus_SUCCESS,
			Name:     new(string),
			UserName: new(string),
			UserId: &userservice.UUID{
				Value: "6c056fb3-7efc-483a-9506-4336456ac79f",
			},
		},
	}
	metaStorage := &metaStorageFake{}
	metaStorage.Store(context.Background(), &SignInMeta{
		SignInKey:   [16]byte{},
		LastRequest: now.Add(-30 * time.Second),
		Phone:       "+7999888776",
		Code:        "123456",
		UserId:      uuid.MustParse("6c056fb3-7efc-483a-9506-4336456ac79f"),
		Username:    "",
	})

	smsSender := smsStub{}
	config := &CodeConfig{
		SendFrequency: 1 * time.Minute,
	}

	sender := NewSignInSendCodeService(config, smsSender, metaStorage, userService)

	t.Run("FrequencyExceeded", func(t *testing.T) {
		// Act
		_, err := sender.SendCode(context.Background(), "+7999888776")

		// Assert
		if assert.Error(t, err) {
			assert.Equal(t, ErrSendCodeFreqExceeded, err)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		defer func(r *userservice.UserResponse) {
			userService.resp = r
		}(userService.resp)

		userService.resp = &userservice.UserResponse{
			Status: userservice.UserResponseStatus_NOT_FOUND,
		}

		// Act
		_, err := sender.SendCode(context.Background(), "+79888888888")

		// Assert
		if assert.Error(t, err) {
			assert.Equal(t, ErrUserNotFound, err)
		}
	})
}

type userServiceMock struct {
	resp *userservice.UserResponse
}

func (s userServiceMock) GetUser(ctx context.Context, in *userservice.UserRequest,
	opts ...grpc.CallOption) (*userservice.UserResponse, error) {
	return s.resp, nil
}

func (s userServiceMock) CreateUser(ctx context.Context, in *userservice.CreateUserRequest,
	opts ...grpc.CallOption) (*userservice.CreateUserResponse, error) {
	panic("why do you use it here?")
}

type metaStorageFake struct {
	s []*SignInMeta
}

func (s *metaStorageFake) FindMetaByPhone(_ context.Context, phone string) (*SignInMeta, bool, error) {
	for _, meta := range s.s {
		if meta.Phone == phone {
			return meta, true, nil
		}
	}
	return nil, false, nil
}

func (s *metaStorageFake) Store(_ context.Context, meta *SignInMeta) error {
	s.s = append(s.s, meta)
	return nil
}

type smsStub struct{}

func (s smsStub) SendSms(_ context.Context, _ string, _ string) (*smsaero_golang.SendSms, error) {
	return nil, nil
}
