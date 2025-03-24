package handlers

import (
	"context"
	"errors"
	"log"
	"regexp"

	pb "github.com/chakchat/chakchat-backend/user-service/internal/grpcservice"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userService services.UserService
}

func NewUserServer(userService services.UserService) *UserServer {
	return &UserServer{userService: userService}
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {

	user, err := s.userService.GetUser(ctx, req.PhoneNumber)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			log.Printf("Failed find user and returns not found, %s", err)
			return &pb.UserResponse{
				Status: pb.UserResponseStatus_NOT_FOUND,
			}, nil
		}
		log.Printf("Unknown fail: %s", err)
		return &pb.UserResponse{
			Status: pb.UserResponseStatus_FAILED,
		}, nil
	}
	log.Println("Pass finding")
	return &pb.UserResponse{
		Status:   pb.UserResponseStatus_SUCCESS,
		Name:     &user.Name,
		UserName: &user.Username,
		UserId:   &pb.UUID{Value: user.ID.String()},
	}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	matchedPhone, _ := regexp.MatchString(`^[79]9\d{9}$`, req.PhoneNumber)
	matchedUsername, _ := regexp.MatchString(`^[a-z][_a-z0-9]{2,19}$`, req.Username)
	matchedName := len(req.Name) <= 50 && len(req.Name) > 0
	if !matchedPhone || !matchedUsername || !matchedName {
		return &pb.CreateUserResponse{
			Status: pb.CreateUserStatus_VALIDATION_FAILED,
		}, nil
	}

	user, err := s.userService.CreateUser(ctx, &models.User{
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.PhoneNumber,
	})

	if err != nil {
		if errors.Is(err, services.ErrAlreadyExists) {
			return &pb.CreateUserResponse{
				Status: pb.CreateUserStatus_ALREADY_EXISTS,
			}, nil
		}
		return &pb.CreateUserResponse{
			Status: pb.CreateUserStatus_CREATE_FAILED,
		}, nil
	}
	return &pb.CreateUserResponse{
		Status:   pb.CreateUserStatus_CREATED,
		UserId:   &pb.UUID{Value: user.ID.String()},
		Name:     &user.Name,
		UserName: &user.Username,
	}, nil
}
