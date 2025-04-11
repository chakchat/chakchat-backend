package services

import (
	"context"
	"time"

	"github.com/chakchat/chakchat-backend/live-connection-service/internal/storage"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/ws"
)

type StatusService struct {
	storage *storage.OnlineStorage
	hub     *ws.Hub
}

type StatusResponse struct {
	UserId     string
	Status     bool
	LastOnline string
}

func NewStatusService(storage *storage.OnlineStorage, hub *ws.Hub) *StatusService {
	return &StatusService{
		storage: storage,
		hub:     hub,
	}
}

func (s *StatusService) GetStatus(ctx context.Context, userIds []string) (map[string]StatusResponse, error) {
	dbStatus, err := s.storage.GetOnlineStatus(ctx, userIds)
	if err != nil {
		return nil, err
	}

	wsStatus := s.hub.GetOnlineStatus(userIds)

	result := make(map[string]StatusResponse)
	for _, id := range userIds {
		result[id] = StatusResponse{
			UserId:     id,
			Status:     wsStatus[id] || time.Since(dbStatus[id].LastOnline) < 10*time.Second,
			LastOnline: dbStatus[id].LastOnline.Format(time.RFC3339),
		}
	}

	return result, nil
}
