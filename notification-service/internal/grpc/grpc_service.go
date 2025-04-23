package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/notification-service/internal/user"
	"github.com/gofrs/uuid"
)

type GRPCClients struct {
	userService user.UserServiceClient
}

func NewGrpcClients(userService user.UserServiceClient) *GRPCClients {
	return &GRPCClients{
		userService: userService,
	}
}

func (c *GRPCClients) GetChatType() (string, error) {
	return "", nil
}
func (c *GRPCClients) GetGroupName() (string, error) {
	return "", nil
}
func (c *GRPCClients) GetName(ctx context.Context, userId uuid.UUID) (string, error) {
	resp, err := c.userService.GetName(ctx, &user.GetNameRequest{
		UserId: userId.String(),
	})

	if err != nil {
		return "", fmt.Errorf("get name gRPC call failed: %s", err)
	}

	switch resp.Status {
	case user.UserResponseStatus_SUCCESS:
		return *resp.Name, nil
	case user.UserResponseStatus_NOT_FOUND:
		return "", errors.New("user not found")
	case user.UserResponseStatus_FAILED:
		return "", errors.New("unknown gRPC GetName() error")
	}

	return *resp.Name, nil
}
