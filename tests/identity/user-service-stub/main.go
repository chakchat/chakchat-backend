package main

import (
	"context"
	"log"
	"net"
	"user-service-stub/userservice"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("tcp listen failed: %s", err)
	}

	ser := grpc.NewServer()
	ser.RegisterService(&userservice.UserService_ServiceDesc, NewServerStub())

	err = ser.Serve(lis)
	if err != nil {
		log.Fatalf("gRPC server failed: %s", err)
	}
}

type ServerStub struct {
	userservice.UnimplementedUserServiceServer
}

func NewServerStub() *ServerStub {
	return &ServerStub{}
}

func (ServerStub) GetUser(ctx context.Context, req *userservice.UserRequest) (*userservice.UserResponse, error) {
	phone := req.GetPhoneNumber()
	// 79********1 phone numbers have existing owners
	if phone[len(phone)-1] == '1' {
		username := "user_with_phone_" + phone
		name := "User with phone " + phone
		id := uuid.Nil.String()
		id = id[:len(id)-11] + phone
		return &userservice.UserResponse{
			Status:   userservice.UserResponseStatus_SUCCESS,
			Name:     &name,
			UserName: &username,
			UserId:   &userservice.UUID{Value: id},
		}, nil
	}

	// 79********2 phone numbers cause fail
	if phone[len(phone)-1] == '2' {
		return &userservice.UserResponse{
			Status: userservice.UserResponseStatus_FAILED,
		}, nil
	}

	// Other phone numbers don't have existing owners
	return &userservice.UserResponse{
		Status: userservice.UserResponseStatus_NOT_FOUND,
	}, nil
}

func (ServerStub) CreateUser(ctx context.Context, req *userservice.CreateUserRequest) (*userservice.CreateUserResponse, error) {
	panic("it is not implemented for now")
}

var _ userservice.UserServiceServer = ServerStub{}

type User struct {
	Id       uuid.UUID
	Phone    string
	Username string
}
