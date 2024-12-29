package services

import (
	"context"
	"testing"
	"time"

	"github.com/chakchat/chakchat/backend/identity/internal/userservice"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_Success(t *testing.T) {
	// Arrange
	userService := usersServiceMock{
		resp: &userservice.UserResponse{
			Status:   userservice.UserResponseStatus_SUCCESS,
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

	sender := NewSendCodeService(&config, smsSender, &metaStorage, userService)

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

	userService := &usersServiceMock{
		resp: &userservice.UserResponse{
			Status:   userservice.UserResponseStatus_SUCCESS,
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

	sender := NewSendCodeService(config, smsSender, metaStorage, userService)

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

type usersServiceMock struct {
	resp *userservice.UserResponse
}

func (s usersServiceMock) GetUser(ctx context.Context, in *userservice.UserRequest,
	opts ...grpc.CallOption) (*userservice.UserResponse, error) {
	return s.resp, nil
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

func (s smsStub) SendSms(_ context.Context, _ string, _ string) error {
	return nil
}
