package proto

import (
	"context"
	"errors"
	"log"

	"github.com/chakchat/chakchat-backend/identity-service/internal/proto/identity"
	"github.com/chakchat/chakchat-backend/identity-service/internal/storage"
	"github.com/google/uuid"
)

type GRPCService struct {
	deviceStorage *storage.DeviceStorage
	identity.UnimplementedIdentityServiceServer
}

func NewGRPCServer(deviceStorage *storage.DeviceStorage) *GRPCService {
	return &GRPCService{
		deviceStorage: deviceStorage,
	}
}

func (s *GRPCService) GetDeviceToken(ctx context.Context, req *identity.DeviceTokenRequest) (*identity.DeviceTokenResponse, error) {

	userId, err := uuid.Parse(req.UserId.Value)
	if err != nil {
		log.Printf("Can't parse userId")
		return &identity.DeviceTokenResponse{
			Status: identity.DeviceTokenResponseStatus_FAILED,
		}, nil
	}

	token, err := s.deviceStorage.GetDeviceTokenByID(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Printf("No device token in data base, %s", err)
			return &identity.DeviceTokenResponse{
				Status: identity.DeviceTokenResponseStatus_NOT_FOUND,
			}, nil
		}
		log.Printf("Unknown fail: %s", err)
		return &identity.DeviceTokenResponse{
			Status: identity.DeviceTokenResponseStatus_FAILED,
		}, nil
	}
	log.Printf("Successfulle get device token")
	return &identity.DeviceTokenResponse{
		Status:      identity.DeviceTokenResponseStatus_SUCCESS,
		DeviceToken: token,
	}, nil
}
