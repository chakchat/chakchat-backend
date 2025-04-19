package storage

import (
	"context"
	"errors"
	"time"

	"github.com/chakchat/chakchat-backend/shared/go/postgres"
	"github.com/lib/pq"
)

var ErrNotFound = errors.New("not found")

type OnlineResponse struct {
	Status     string
	LastOnline time.Time
}

type OnlineStorage struct {
	db postgres.SQLer
}

func NewOnlineStorage(db postgres.SQLer) *OnlineStorage {
	return &OnlineStorage{db: db}
}

func (s *OnlineStorage) UpdateLastPing(ctx context.Context, userId string) error {
	query := `
		INSERT INTO user_online_status (user_id, last_ping)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET last_ping = $2
	`
	_, err := s.db.Exec(ctx, query, userId, time.Now())
	return err
}

func (s *OnlineStorage) GetOnlineStatus(ctx context.Context, userIds []string) (map[string]OnlineResponse, error) {
	query := `
		SELECT user_id, last_ping
		FROM user_online_status
		WHERE user_id = ANY($1)
	`

	rows, err := s.db.Query(ctx, query, pq.Array(userIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]OnlineResponse)
	now := time.Now()

	for rows.Next() {
		var userID string
		var lastPing time.Time
		if err := rows.Scan(&userID, &lastPing); err != nil {
			continue
		}

		var status string
		if now.Sub(lastPing) < 10*time.Second {
			status = "online"
		} else {
			status = "offline"
		}
		result[userID] = OnlineResponse{
			Status:     status,
			LastOnline: lastPing,
		}
	}

	for _, id := range userIds {
		if _, exists := result[id]; !exists {
			result[id] = OnlineResponse{Status: "offline"}
		}
	}

	return result, nil
}
