package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secpersonal"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SecretPersonalChatRepository struct {
	db *pgx.Conn
}

func NewSecretPersonalChatRepository(db *pgx.Conn) *SecretPersonalChatRepository {
	return &SecretPersonalChatRepository{
		db: db,
	}
}

func (r *SecretPersonalChatRepository) FindById(ctx context.Context, id domain.ChatID) (*secpersonal.SecretPersonalChat, error) {
	q := `
	SELECT 
		c.chat_id, 
		c.created_at, 
		p.expiration_seconds,
		(SELECT ARRAY_AGG(m.user_id) FROM messaging.membership m WHERE m.chat_id = c.chat_id), 
	FROM messaging.chat c
		JOIN messaging.personal_chat p ON p.chat_id = c.chat_id
	WHERE c.chat_id = $1`

	row := r.db.QueryRow(ctx, q, id)

	var (
		chatID            uuid.UUID
		createdAt         time.Time
		expirationSeconds *int
		members           []uuid.UUID
	)

	err := row.Scan(&chatID, &createdAt, &expirationSeconds, &members)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("getting secret personal chat failed: %s", err)
	}

	exp := (*time.Duration)(nil)
	if expirationSeconds != nil {
		cp := time.Duration(*expirationSeconds) * time.Second
		exp = &cp
	}

	return &secpersonal.SecretPersonalChat{
		SecretChat: domain.SecretChat{
			Chat: domain.Chat{
				ID:        domain.ChatID(chatID),
				CreatedAt: domain.Timestamp(createdAt.Unix()),
			},
			Exp: exp,
		},
		Members: [2]domain.UserID(userIDs(members)),
	}, nil
}

func (r *SecretPersonalChatRepository) FindByMembers(ctx context.Context, members [2]domain.UserID) (*secpersonal.SecretPersonalChat, error) {
	q := `
	SELECT 
		c.chat_id, 
		MAX(c.created_at),
		MAX(sp.expiration_seconds),
	FROM messaging.membership m
		JOIN messaging.chat c ON c.chat_id = m.chat_id
		JOIN messaging.secret_personal_chat sp ON sp.chat_id = c.chat_id
	WHERE m.user_id = $1 OR m.user_id = $2
	GROUP BY c.chat_id
	HAVING COUNT(DISTINCT m.user_id) = 2`

	row := r.db.QueryRow(ctx, q, members[0], members[1])

	var (
		chatID            uuid.UUID
		createdAt         time.Time
		expirationSeconds *int
	)

	err := row.Scan(&chatID, &createdAt, &expirationSeconds)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("getting secret personal chat failed: %s", err)
	}

	exp := (*time.Duration)(nil)
	if expirationSeconds != nil {
		cp := time.Duration(*expirationSeconds) * time.Second
		exp = &cp
	}

	return &secpersonal.SecretPersonalChat{
		SecretChat: domain.SecretChat{
			Chat: domain.Chat{
				ID:        domain.ChatID(chatID),
				CreatedAt: domain.Timestamp(createdAt.Unix()),
			},
			Exp: exp,
		},
		Members: members,
	}, nil
}

func (r *SecretPersonalChatRepository) Update(ctx context.Context, chat *secpersonal.SecretPersonalChat) (*secpersonal.SecretPersonalChat, error) {
	q := `
	UPDATE messaging.secret_personal_chat
	SET expiration_seconds = $2
	WHERE chat_id = $1`

	_, err := r.db.Exec(ctx, q, chat.ID, int(chat.Expiration().Seconds()))
	return chat, err
}
