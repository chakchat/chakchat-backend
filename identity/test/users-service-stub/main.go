package main

import (
	"context"
	"log"
	"net"
	"user-service-stub/userservice"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var existingUser = &User{
	Id:       uuid.MustParse("11111111-1111-1111-1111-111111111111"),
	Phone:    "79111111111",
	Username: "bob",
}

var erroringUser = &User{
	Id:       uuid.MustParse("22222222-2222-2222-2222-222222222222"),
	Phone:    "79222222222",
	Username: "bib",
}

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("tcp listen failed: %s", err)
	}

	ser := grpc.NewServer()
	ser.RegisterService(&userservice.UsersService_ServiceDesc, NewServerStub())

	err = ser.Serve(lis)
	if err != nil {
		log.Fatalf("gRPC server failed: %s", err)
	}
}

type ServerStub struct {
	userservice.UnimplementedUsersServiceServer
}

func NewServerStub() *ServerStub {
	return &ServerStub{}
}

func (s ServerStub) GetUser(ctx context.Context, req *userservice.UserRequest) (*userservice.UserResponse, error) {
	if req.GetPhoneNumber() == existingUser.Phone {
		return &userservice.UserResponse{
			Status:   userservice.UserResponseStatus_SUCCESS,
			UserName: &existingUser.Username,
			UserId: &userservice.UUID{
				Value: existingUser.Id.String(),
			},
		}, nil
	}

	if req.GetPhoneNumber() == erroringUser.Phone {
		return &userservice.UserResponse{
			Status: userservice.UserResponseStatus_FAILED,
		}, nil
	}

	return &userservice.UserResponse{
		Status: userservice.UserResponseStatus_NOT_FOUND,
	}, nil
}

var _ userservice.UsersServiceServer = ServerStub{}

type User struct {
	Id       uuid.UUID
	Phone    string
	Username string
}
