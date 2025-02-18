package handlers

import (
	"context"
	"errors"
	"regexp"

	pb "github.com/chakchat/chakchat-backend/user-service/internal/grpcservice"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userService services.UserService
}

func NewUserServer(userService services.UserService) *UserServer {
	return &UserServer{userService: userService}
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	phone := req.PhoneNumber

	user, err := s.userService.GetUser(ctx, phone)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return &pb.UserResponse{
				Status: pb.UserResponseStatus_NOT_FOUND,
			}, nil
		}
		return &pb.UserResponse{
			Status: pb.UserResponseStatus_FAILED,
		}, nil
	}
	return &pb.UserResponse{
		Status:   pb.UserResponseStatus_SUCCESS,
		Name:     &user.Name,
		UserName: &user.Username,
		UserId:   &pb.UUID{Value: user.ID.String()},
	}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	name, username, phone := req.Name, req.Username, req.PhoneNumber
	matchedPhone, _ := regexp.MatchString(`^\+79\d{9}$`, phone)
	matchedUsername, _ := regexp.MatchString(`^[a-z][_a-z0-9]{2,19}$`, username)
	matchedName := len(name) <= 50 && len(name) > 0
	if !matchedPhone || !matchedUsername || !matchedName {
		return &pb.CreateUserResponse{
			Status: pb.CreateUserStatus_VALIDATION_FAILED,
		}, nil
	}

	user, err := s.userService.CreateUser(ctx, &storage.User{
		Name:     name,
		Username: username,
		Phone:    phone,
	})

	if err != nil {
		if errors.Is(err, services.ErrAlreadyExists) {
			return &pb.CreateUserResponse{
				Status: pb.CreateUserStatus_ALREADY_EXISTS,
			}, nil
		} else {
			return &pb.CreateUserResponse{
				Status: pb.CreateUserStatus_CREATE_FAILED,
			}, nil
		}
	}
	return &pb.CreateUserResponse{
		Status:   pb.CreateUserStatus_CREATED,
		UserId:   &pb.UUID{Value: user.ID.String()},
		Name:     &user.Name,
		UserName: &user.Username,
	}, nil
}
