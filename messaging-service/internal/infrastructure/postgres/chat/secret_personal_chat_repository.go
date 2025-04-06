package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secpersonal"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SecretPersonalChatRepository struct{}

func NewSecretPersonalChatRepository() *SecretPersonalChatRepository {
	return &SecretPersonalChatRepository{}
}

func (r *SecretPersonalChatRepository) FindById(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID,
) (*secpersonal.SecretPersonalChat, error) {
	q := `
	SELECT 
		c.chat_id, 
		c.created_at, 
		p.expiration_seconds,
		(SELECT ARRAY_AGG(m.user_id) FROM messaging.membership m WHERE m.chat_id = c.chat_id)
	FROM messaging.chat c
		JOIN messaging.personal_chat p ON p.chat_id = c.chat_id
	WHERE c.chat_id = $1`

	row := db.QueryRow(ctx, q, id)

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

func (r *SecretPersonalChatRepository) FindByMembers(
	ctx context.Context, db storage.ExecQuerier, members [2]domain.UserID,
) (*secpersonal.SecretPersonalChat, error) {
	q := `
	SELECT 
		c.chat_id, 
		MAX(c.created_at),
		MAX(sp.expiration_seconds)
	FROM messaging.membership m
		JOIN messaging.chat c ON c.chat_id = m.chat_id
		JOIN messaging.secret_personal_chat sp ON sp.chat_id = c.chat_id
	WHERE m.user_id = $1 OR m.user_id = $2
	GROUP BY c.chat_id
	HAVING COUNT(DISTINCT m.user_id) = 2`

	row := db.QueryRow(ctx, q, members[0], members[1])

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

func (r *SecretPersonalChatRepository) Update(
	ctx context.Context, db storage.ExecQuerier, chat *secpersonal.SecretPersonalChat,
) (*secpersonal.SecretPersonalChat, error) {
	q := `
	UPDATE messaging.secret_personal_chat
	SET expiration_seconds = $2
	WHERE chat_id = $1`

	var exp *int
	if chat.Exp != nil {
		cp := int(chat.Exp.Seconds())
		exp = &cp
	}

	_, err := db.Exec(ctx, q, chat.ID, exp)
	return chat, err
}

func (r *SecretPersonalChatRepository) Create(
	ctx context.Context, db storage.ExecQuerier, chat *secpersonal.SecretPersonalChat,
) (*secpersonal.SecretPersonalChat, error) {
	{
		q := `
		INSERT INTO messaging.chat
		(chat_id, chat_type, created_at)
		VALUES ($1, 'personal', $2)`

		now := time.Now()
		_, err := db.Exec(ctx, q, chat.ID, now)
		if err != nil {
			return nil, err
		}
		chat.CreatedAt = domain.Timestamp(now.Unix())
	}
	{
		var exp *int
		if chat.Exp != nil {
			cp := int(chat.Exp.Seconds())
			exp = &cp
		}
		q := `INSERT INTO messaging.personal_chat (chat_id, expiration_seconds) VALUES ($1, $2)`
		_, err := db.Exec(ctx, q, chat.ID, exp)
		if err != nil {
			return nil, err
		}
	}
	{
		q := `INSERT INTO messaging.membership (chat_id, user_id) VALUES ($1, $2), ($1, $3)`
		_, err := db.Exec(ctx, q, chat.ID, chat.Members[0], chat.Members[1])
		if err != nil {
			return nil, err
		}
	}
	return chat, nil
}

func (r *SecretPersonalChatRepository) Delete(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID,
) error {
	q := `DELETE FROM messaging.chat WHERE chat_id = $1`
	_, err := db.Exec(ctx, q, id)
	return err
}
