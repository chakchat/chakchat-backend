package main

import (
	"context"
	"os"
	"test/userservice"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGrpc(t *testing.T) {
	addr := os.Getenv("USER_SERVICE_ADDR")
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Connecting to UserService failed: %s", err)
	}
	client := userservice.NewUserServiceClient(conn)

	{
		resp, err := client.CreateUser(context.Background(), &userservice.CreateUserRequest{
			PhoneNumber: "79012345678",
			Username:    "user1",
			Name:        "Some user 1",
		})
		require.NoError(t, err, "gRPC call failed")
		require.Equal(t, userservice.CreateUserStatus_CREATED, resp.Status)

		require.NotNil(t, resp.Name)
		require.NotNil(t, resp.UserName)
		require.NotNil(t, resp.UserId)

		require.Equal(t, "user1", *resp.UserName)
		require.Equal(t, "Some user 1", *resp.Name)

		id, err := uuid.Parse(resp.UserId.Value)
		require.NoError(t, err)
		require.NotZero(t, id)
	}
	{
		resp, err := client.CreateUser(context.Background(), &userservice.CreateUserRequest{
			PhoneNumber: "79012345678",
			Username:    "user1",
			Name:        "Some user 1",
		})
		require.NoError(t, err, "gRPC call failed")
		require.Equal(t, userservice.CreateUserStatus_ALREADY_EXISTS, resp.Status)
	}
	{
		resp, err := client.CreateUser(context.Background(), &userservice.CreateUserRequest{
			PhoneNumber: "79012345",
			Username:    "user123",
			Name:        "Some user 1",
		})
		require.NoError(t, err, "gRPC call failed")
		require.Equal(t, userservice.CreateUserStatus_VALIDATION_FAILED, resp.Status)
	}
	{
		resp, err := client.CreateUser(context.Background(), &userservice.CreateUserRequest{
			PhoneNumber: "79012345622",
			Username:    "_user2",
			Name:        "Some user 1",
		})
		require.NoError(t, err, "gRPC call failed")
		require.Equal(t, userservice.CreateUserStatus_VALIDATION_FAILED, resp.Status)
	}
	{
		resp, err := client.CreateUser(context.Background(), &userservice.CreateUserRequest{
			PhoneNumber: "79012345622",
			Username:    "user3",
			Name:        "TOO LONG ajjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj",
		})
		require.NoError(t, err, "gRPC call failed")
		require.Equal(t, userservice.CreateUserStatus_VALIDATION_FAILED, resp.Status)
	}
	{
		resp, err := client.GetUser(context.Background(), &userservice.UserRequest{
			PhoneNumber: "790123456789",
		})
		require.NoError(t, err, "gRPC call failed")

		require.Equal(t, userservice.UserResponseStatus_SUCCESS, resp.Status)

		require.NotNil(t, resp.Name)
		require.NotNil(t, resp.UserName)
		require.NotNil(t, resp.UserId)

		require.Equal(t, "user1", *resp.UserName)
		require.Equal(t, "Some user 1", *resp.Name)

		id, err := uuid.Parse(resp.UserId.Value)
		require.NoError(t, err)
		require.NotZero(t, id)
	}
	{
		resp, err := client.GetUser(context.Background(), &userservice.UserRequest{
			PhoneNumber: "79000000000",
		})
		require.NoError(t, err, "gRPC call failed")

		require.Equal(t, userservice.UserResponseStatus_NOT_FOUND, resp.Status)
	}
}
