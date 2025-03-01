package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/personal"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func userIDs(arr pgtype.Array[uuid.UUID]) []domain.UserID {
	res := make([]domain.UserID, arr.Dims[0].Length)
	for i := range res {
		res[i] = domain.UserID(arr.Index(i).(uuid.UUID))
	}
	return res
}

type PersonalChatRepository struct {
	db *pgx.Conn
}

func NewPersonalChatRepository(db *pgx.Conn) *PersonalChatRepository {
	return &PersonalChatRepository{
		db: db,
	}
}

func (r *PersonalChatRepository) FindById(ctx context.Context, id domain.ChatID) (*personal.PersonalChat, error) {
	q := `
	SELECT 
		c.chat_id, 
		c.created_at, 
		(SELECT ARRAY_AGG(user_id) FROM messaging.membership m WHERE m.chat_id = c.chat_id), 
		(SELECT ARRAY_AGG(user_id) FROM messaging.blockings b WHERE b.chat_id = c.chat_id)
	FROM messaging.chat c
		JOIN messaging.personal_chat p ON p.chat_id = c.chat_id
	WHERE c.chat_id = $1`

	row := r.db.QueryRow(ctx, q, id)

	var (
		chatID    uuid.UUID
		createdAt time.Time
		members   pgtype.Array[uuid.UUID]
		blockedBy pgtype.Array[uuid.UUID]
	)

	err := row.Scan(&chatID, &createdAt, &members, &blockedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("getting personal chat failed: %s", err)
	}

	return &personal.PersonalChat{
		Chat: domain.Chat{
			ID:        id,
			CreatedAt: domain.Timestamp(createdAt.Unix()),
		},
		Members:   [2]domain.UserID(userIDs(members)),
		BlockedBy: userIDs(blockedBy),
	}, nil
}
