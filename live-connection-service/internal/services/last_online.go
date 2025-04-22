package services

import (
	"context"
	"time"

	"github.com/chakchat/chakchat-backend/live-connection-service/internal/storage"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/ws"
	"github.com/google/uuid"
)

type StatusService struct {
	storage *storage.OnlineStorage
	hub     *ws.Hub
}

type StatusResponse struct {
	UserId     uuid.UUID
	Status     bool
	LastOnline string
}

func NewStatusService(storage *storage.OnlineStorage, hub *ws.Hub) *StatusService {
	return &StatusService{
		storage: storage,
		hub:     hub,
	}
}

func (s *StatusService) GetStatus(ctx context.Context, userIds []uuid.UUID) (map[uuid.UUID]StatusResponse, error) {
	dbStatus, err := s.storage.GetOnlineStatus(ctx, userIds)
	if err != nil {
		return nil, err
	}

	wsStatus := s.hub.GetOnlineStatus(userIds)

	result := make(map[uuid.UUID]StatusResponse)
	for _, id := range userIds {
		result[id] = StatusResponse{
			UserId:     id,
			Status:     wsStatus[id] || time.Since(dbStatus[id].LastOnline) < 10*time.Second,
			LastOnline: dbStatus[id].LastOnline.Format(time.RFC3339),
		}
	}

	return result, nil
}
