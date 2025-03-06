package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/personal"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		(SELECT ARRAY_AGG(m.user_id) FROM messaging.membership m WHERE m.chat_id = c.chat_id), 
		(SELECT ARRAY_AGG(b.user_id) FROM messaging.blocking b WHERE b.chat_id = c.chat_id)
	FROM messaging.chat c
		JOIN messaging.personal_chat p ON p.chat_id = c.chat_id
	WHERE c.chat_id = $1`

	row := r.db.QueryRow(ctx, q, id)

	var (
		chatID    uuid.UUID
		createdAt time.Time
		members   []uuid.UUID
		blockedBy []uuid.UUID
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
		(SELECT ARRAY_AGG(user_id) FROM messaging.blocking b WHERE b.chat_id = m.chat_id)
	FROM messaging.membership m
		JOIN messaging.chat c ON c.chat_id = m.chat_id
		JOIN messaging.personal_chat p ON p.chat_id = m.chat_id
	WHERE m.user_id = $1 OR m.user_id = $2
	GROUP BY m.chat_id
	HAVING COUNT(DISTINCT m.user_id) = 2
	`

	row := r.db.QueryRow(ctx, q, uuid.UUID(members[0]), uuid.UUID(members[1]))

	var (
		chatID    uuid.UUID
		createdAt time.Time
		blockedBy []uuid.UUID
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

func (r *PersonalChatRepository) Update(ctx context.Context, chat *personal.PersonalChat) (*personal.PersonalChat, error) {
	blockings, err := r.getBlockings(ctx, uuid.UUID(chat.ID))
	if err != nil {
		return nil, err
	}

	toAdd := sliceMisses(chat.BlockedBy, blockings)
	if len(toAdd) != 0 {
		if err := r.addBlocking(ctx, uuid.UUID(chat.ID), uuids(toAdd)); err != nil {
			return nil, err
		}
	}

	toDelete := sliceMisses(blockings, chat.BlockedBy)
	if len(toDelete) != 0 {
		if err := r.deleteBlocking(ctx, uuid.UUID(chat.ID), uuids(toDelete)); err != nil {
			return nil, err
		}
	}

	return chat, nil
}

func (r *PersonalChatRepository) Create(ctx context.Context, chat *personal.PersonalChat) (*personal.PersonalChat, error) {
	{
		q := `
		INSERT INTO messaging.chat
		(chat_id, chat_type, created_at)
		VALUES ($1, 'personal', $2)
		`
		now := time.Now()
		_, err := r.db.Exec(ctx, q, chat.ID, now)
		if err != nil {
			return nil, err
		}
		chat.CreatedAt = domain.Timestamp(now.Unix())
	}
	{
		q := `INSERT INTO messaging.personal_chat (chat_id) VALUES ($1)`
		_, err := r.db.Exec(ctx, q, chat.ID)
		if err != nil {
			return nil, err
		}
	}
	{
		q := `INSERT INTO messaging.membership (chat_id, user_id) VALUES ($1, $2), ($1, $3)`
		_, err := r.db.Exec(ctx, q, chat.ID, chat.Members[0], chat.Members[1])
		if err != nil {
			return nil, err
		}
	}
	return chat, nil
}

func (r *PersonalChatRepository) Delete(ctx context.Context, id domain.ChatID) error {
	q := `DELETE FROM messaging.chat WHERE chat_id = $1`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

func (r *PersonalChatRepository) addBlocking(ctx context.Context, chatId uuid.UUID, toAdd []uuid.UUID) error {
	q := `INSERT INTO messaging.blocking (chat_id, user_id) VALUES `

	valExprs := make([]string, 0, len(toAdd))
	args := make([]any, 0, 2*len(toAdd))
	for _, userId := range toAdd {
		argI := len(args) + 1
		valExprs = append(valExprs, fmt.Sprintf("($%d, $%d)", argI, argI+1))
		args = append(args, chatId, userId)
	}

	q += strings.Join(valExprs, ", ")

	_, err := r.db.Exec(ctx, q, args...)
	return err
}

func (r *PersonalChatRepository) deleteBlocking(ctx context.Context, chatId uuid.UUID, toDelete []uuid.UUID) error {
	q := `DELETE FROM messaging.blocking WHERE chat_id = $1 AND user_id = ANY($2)`

	_, err := r.db.Exec(ctx, q, chatId, toDelete)
	return err
}

func (r *PersonalChatRepository) getBlockings(ctx context.Context, id uuid.UUID) ([]domain.UserID, error) {
	q := `SELECT user_id FROM messaging.blocking WHERE chat_id = $1`

	rows, err := r.db.Query(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("get blockings query failed: %s", err)
	}

	res := make([]domain.UserID, 0)
	for rows.Next() {
		var curr domain.UserID
		if err := rows.Scan(&curr); err != nil {
			return nil, fmt.Errorf("scanning rows failed: %s", err)
		}
		res = append(res, curr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sql rows returned an error: %s", err)
	}

	return res, nil
}
