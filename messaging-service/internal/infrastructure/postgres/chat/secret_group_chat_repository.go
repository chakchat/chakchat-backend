package chat

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secgroup"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SecretGroupChatRepository struct{}

func NewSecretGroupChatRepository() *SecretGroupChatRepository {
	return &SecretGroupChatRepository{}
}

func (r *SecretGroupChatRepository) FindById(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID,
) (*secgroup.SecretGroupChat, error) {
	q := `
	SELECT c.chat_id, c.created_at, g.admin_id, g.group_name, g.group_photo, g.group_description, g.expiration_seconds,
		(SELECT ARRAY_AGG(m.user_id) FROM messaging.membership m WHERE m.chat_id = c.chat_id)
	FROM messaging.chat c
		JOIN messaging.secret_group_chat g ON g.chat_id = c.chat_id
	WHERE c.chat_id = $1`

	row := db.QueryRow(ctx, q, id)

	var (
		chatID            uuid.UUID
		createdAt         time.Time
		adminID           uuid.UUID
		name              string
		photo             string
		description       string
		expirationSeconds *int
		members           []uuid.UUID
	)
	err := row.Scan(&chatID, &createdAt, &adminID, &name, &photo, &description, &expirationSeconds, &members)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("getting group chat failed: %s", err)
	}

	var exp *time.Duration
	if expirationSeconds != nil {
		cp := time.Duration(*expirationSeconds) * time.Second
		exp = &cp
	}
	return &secgroup.SecretGroupChat{
		SecretChat: domain.SecretChat{
			Chat: domain.Chat{
				ID:        id,
				CreatedAt: domain.Timestamp(createdAt.Unix()),
			},
			Exp: exp,
		},
		Admin:       domain.UserID(adminID),
		Members:     userIDs(members),
		Name:        name,
		Description: description,
		GroupPhoto:  domain.URL(photo),
	}, nil
}

func (r *SecretGroupChatRepository) Update(
	ctx context.Context, db storage.ExecQuerier, g *secgroup.SecretGroupChat,
) (*secgroup.SecretGroupChat, error) {
	members, err := r.getMembers(ctx, db, g.ID)
	if err != nil {
		return nil, err
	}

	toAdd := sliceMisses(g.Members, members)
	if len(toAdd) != 0 {
		if err := r.addMembers(ctx, db, g.ID, toAdd); err != nil {
			return nil, err
		}
	}

	toDelete := sliceMisses(members, g.Members)
	if len(toDelete) != 0 {
		if err := r.deleteMembers(ctx, db, uuid.UUID(g.ID), uuids(toDelete)); err != nil {
			return nil, err
		}
	}

	q := `
	UPDATE messaging.secret_group_chat 
	SET admin_id = $2, 
		group_name = $3, 
		group_photo = $4, 
		group_description = $5
		expiration_seconds = $6
	WHERE chat_id = $1`

	var exp *int
	if g.Exp != nil {
		cp := int(g.Exp.Seconds())
		exp = &cp
	}

	_, err = db.Exec(ctx, q, g.ID, g.Admin, g.Name, g.GroupPhoto, g.Description, exp)
	if err != nil {
		return nil, fmt.Errorf("updating group chat failed: %s", err)
	}

	return g, err
}

func (r *SecretGroupChatRepository) Create(
	ctx context.Context, db storage.ExecQuerier, g *secgroup.SecretGroupChat,
) (*secgroup.SecretGroupChat, error) {
	{
		q := `
		INSERT INTO messaging.chat
		(chat_id, chat_type, created_at)
		VALUES ($1, 'group', $2)`

		now := time.Now()
		_, err := db.Exec(ctx, q, g.ID, now)
		if err != nil {
			return nil, err
		}
		g.CreatedAt = domain.Timestamp(now.Unix())
	}
	{
		var exp *int
		if g.Exp != nil {
			cp := int(g.Exp.Seconds())
			exp = &cp
		}
		q := `
		INSERT INTO messaging.secret_group_chat
		(chat_id, admin_id, group_name, group_photo, group_description, expiration_seconds)
		VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := db.Exec(ctx, q, g.ID, g.Name, g.Admin, g.GroupPhoto, g.Description, exp)
		if err != nil {
			return nil, err
		}
	}
	{
		q := `INSERT INTO messaging.membership (chat_id, user_id) VALUES `

		valExprs := make([]string, 0, len(g.Members))
		args := make([]any, 0, 2*len(g.Members))
		for _, userId := range g.Members {
			argI := len(args) + 1
			valExprs = append(valExprs, fmt.Sprintf("($%d, $%d)", argI, argI+1))
			args = append(args, g.ID, userId)
		}

		q += strings.Join(valExprs, ", ")

		_, err := db.Exec(ctx, q, args...)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (r *SecretGroupChatRepository) Delete(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID,
) error {
	q := `DELETE FROM messaging.chat WHERE chat_id = $1`
	_, err := db.Exec(ctx, q, id)
	return err
}

func (r *SecretGroupChatRepository) addMembers(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID, toAdd []domain.UserID,
) error {
	q := `INSERT INTO messaging.membership (chat_id, user_id) VALUES `

	valExprs := make([]string, 0, len(toAdd))
	args := make([]any, 0, 2*len(toAdd))
	for _, userId := range toAdd {
		argI := len(args) + 1
		valExprs = append(valExprs, fmt.Sprintf("($%d, $%d)", argI, argI+1))
		args = append(args, id, userId)
	}

	q += strings.Join(valExprs, ", ")

	_, err := db.Exec(ctx, q, args...)
	return err
}

func (r *SecretGroupChatRepository) deleteMembers(
	ctx context.Context, db storage.ExecQuerier, id uuid.UUID, toDelete []uuid.UUID,
) error {
	q := `DELETE FROM messaging.membership WHERE chat_id = $1 AND user_id = ANY($2)`

	_, err := db.Exec(ctx, q, id, toDelete)
	return err
}

func (r *SecretGroupChatRepository) getMembers(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID,
) ([]domain.UserID, error) {
	q := `SELECT user_id FROM messaging.membership WHERE chat_id = $1`

	rows, err := db.Query(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("get members query failed: %s", err)
	}
	defer rows.Close()

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
