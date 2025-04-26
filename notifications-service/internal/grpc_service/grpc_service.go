package grpc_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/notification-service/internal/identity"
	"github.com/chakchat/chakchat-backend/notification-service/internal/user"
	"github.com/gofrs/uuid"
)

var ErrNotFound = errors.New("not found")

type GRPCClients struct {
	userService     user.UserServiceClient
	identityService identity.IdentityServiceClient
}

func NewGrpcClients(userService user.UserServiceClient, identityService identity.IdentityServiceClient) *GRPCClients {
	return &GRPCClients{
		userService:     userService,
		identityService: identityService,
	}
}

func (c *GRPCClients) GetChatType() (string, error) {
	return "", nil
}
func (c *GRPCClients) GetGroupName() (string, error) {
	return "", nil
}
func (c *GRPCClients) GetName(ctx context.Context, userId uuid.UUID) (*string, error) {
	resp, err := c.userService.GetName(ctx, &user.GetNameRequest{
		UserId: userId.String(),
	})

	if err != nil {
		return nil, fmt.Errorf("get name gRPC call failed: %s", err)
	}

	switch resp.Status {
	case user.UserResponseStatus_SUCCESS:
		return resp.Name, nil
	case user.UserResponseStatus_NOT_FOUND:
		return nil, ErrNotFound
	case user.UserResponseStatus_FAILED:
		return nil, errors.New("unknown gRPC GetName() error")
	}

	return resp.Name, nil
}

func (c *GRPCClients) GetDeviceToken(ctx context.Context, userId uuid.UUID) (*string, error) {
	resp, err := c.identityService.GetDeviceTokens(ctx, &identity.DeviceTokenRequest{
		UserId: &identity.UUID{Value: userId.String()},
	})

	if err != nil {
		return nil, fmt.Errorf("get device token gRPC call failed: %s", err)
	}

	switch resp.Status {
	case identity.DeviceTokenResponseStatus_FAILED:
		return nil, errors.New("unknown gRPC GetDeviceToken() error")
	case identity.DeviceTokenResponseStatus_NOT_FOUND:
		return nil, nil
	case identity.DeviceTokenResponseStatus_SUCCESS:
		return resp.DeviceToken, nil
	}
	return resp.DeviceToken, nil
}
