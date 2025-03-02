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
			ID:        domain.ChatID(chatID),
			CreatedAt: domain.Timestamp(createdAt.Unix()),
		},
		Members:   [2]domain.UserID(userIDs(members)),
		BlockedBy: userIDs(blockedBy),
	}, nil
}

func (r *PersonalChatRepository) FindByMembers(ctx context.Context, members [2]domain.UserID) (*personal.PersonalChat, error) {
	q := `
	SELECT 
		m.chat_id, 
		MAX(c.created_at),
		(SELECT ARRAY_AGG(user_id) FROM messaging.blockings b WHERE b.chat_id = c.chat_id)
	FROM messaging.membership m
		JOIN messaging.chat c ON c.chat_id = m.chat_id
		JOIN messaging.personal_chat p ON p.chat_id = m.chat_id
	WHERE m.user_id = $1 OR m.user_id = $2
	GROUP BY m.chat_id
	HAVING COUNT(DISCTINCT m.user_id) = 2
	`

	row := r.db.QueryRow(ctx, q, uuid.UUID(members[0]), uuid.UUID(members[1]))

	var (
		chatID    uuid.UUID
		createdAt time.Time
		blockedBy pgtype.Array[uuid.UUID]
	)

	err := row.Scan(&chatID, &createdAt, &blockedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("getting personal chat failed: %s", err)
	}

	return &personal.PersonalChat{
		Chat: domain.Chat{
			ID:        domain.ChatID(chatID),
			CreatedAt: domain.Timestamp(createdAt.Unix()),
		},
		Members:   members,
		BlockedBy: userIDs(blockedBy),
	}, nil
}
