package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GenericChatRepository struct {
	db *pgx.Conn
}

func NewGenericChatRepository(db *pgx.Conn) *GenericChatRepository {
	return &GenericChatRepository{
		db: db,
	}
}

func (r *GenericChatRepository) GetByMemberID(ctx context.Context, memberID domain.UserID) ([]services.GenericChat, error) {
	// I guess this query is so fucking slow in real production,
	// but for now we have more microservices then users, so nobody cares.
	//
	// TODO: optimize it.
	q := `
	SELECT 
		c.chat_id,
		c.chat_type,
		c.created_at,
		(SELECT ARRAY_AGG(member_id) FROM messaging.membership WHERE chat_id = c.chat_id),
		CASE WHEN c.chat_type = 'personal' THEN (SELECT ARRAY_AGG(b.user_id) FROM messaging.blocking b WHERE b.chat_id = c.chat_id)
			ELSE NULL
		END,
		COALESCE(group_chat.admin_id, secret_group_chat.admin_id),
		COALESCE(group_chat.group_name, secret_group_chat.group_name),
		COALESCE(group_chat.group_photo, secret_group_chat.group_photo),
		COALESCE(group_chat.group_description, secret_group_chat.group_description),
		COALESCE(secret_personal_chat.expiration_seconds, secret_group_chat.expiration_seconds)
	FROM messaging.membership m
		JOIN messaging.chat c ON c.chat_id = m.chat_id
		LEFT JOIN messaging.personal_chat ON personal_chat.chat_id = c.chat_id
		LEFT JOIN messaging.group_chat ON group_chat.chat_id = c.chat_id
		LEFT JOIN messaging.secret_personal_chat ON secret_personal_chat.chat_id = c.chat_id
		LEFT JOIN messaging.secret_group_chat ON secret_group_chat.chat_id = c.chat_id
	WHERE m.member_id = $1`

	rows, err := r.db.Query(ctx, q, memberID)
	if err != nil {
		return nil, fmt.Errorf("getting chats by memberID failed: %s", err)
	}

	res := make([]services.GenericChat, 0)
	for rows.Next() {
		var (
			chatID            uuid.UUID
			chatType          string
			createdAt         time.Time
			members           []uuid.UUID
			blockedBy         []uuid.UUID
			adminID           uuid.UUID
			groupName         string
			groupPhoto        string
			groupDescription  string
			expirationSeconds *int
		)
		err := rows.Scan(&chatID, &chatType, &createdAt, &members,
			&adminID, &groupName, &groupPhoto, &groupDescription, &expirationSeconds)
		if err != nil {
			return nil, err
		}

		res = append(res, r.buildGenericChat(chatID, chatType, createdAt, members, blockedBy,
			adminID, groupName, groupPhoto, groupDescription, expirationSeconds))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *GenericChatRepository) GetByChatID(ctx context.Context, id domain.ChatID) (*services.GenericChat, error) {
	// I guess this query is so fucking slow in real production,
	// but for now we have more microservices then users, so nobody cares.
	//
	// TODO: optimize it.
	q := `
	SELECT 
		c.chat_id,
		c.chat_type,
		c.created_at,
		(SELECT ARRAY_AGG(member_id) FROM messaging.membership WHERE chat_id = c.chat_id),
		CASE WHEN c.chat_type = 'personal' THEN (SELECT ARRAY_AGG(b.user_id) FROM messaging.blocking b WHERE b.chat_id = c.chat_id)
			ELSE NULL
		END,
		COALESCE(group_chat.admin_id, secret_group_chat.admin_id),
		COALESCE(group_chat.group_name, secret_group_chat.group_name),
		COALESCE(group_chat.group_photo, secret_group_chat.group_photo),
		COALESCE(group_chat.group_description, secret_group_chat.group_description),
		COALESCE(secret_personal_chat.expiration_seconds, secret_group_chat.expiration_seconds)
	FROM messaging.chat c
		LEFT JOIN messaging.personal_chat ON personal_chat.chat_id = c.chat_id
		LEFT JOIN messaging.group_chat ON group_chat.chat_id = c.chat_id
		LEFT JOIN messaging.secret_personal_chat ON secret_personal_chat.chat_id = c.chat_id
		LEFT JOIN messaging.secret_group_chat ON secret_group_chat.chat_id = c.chat_id
	WHERE m.chat_id = $1`

	row := r.db.QueryRow(ctx, q, id)

	var (
		chatID            uuid.UUID
		chatType          string
		createdAt         time.Time
		members           []uuid.UUID
		blockedBy         []uuid.UUID
		adminID           uuid.UUID
		groupName         string
		groupPhoto        string
		groupDescription  string
		expirationSeconds *int
	)
	err := row.Scan(&chatID, &chatType, &createdAt, &members, &blockedBy,
		&adminID, &groupName, &groupPhoto, &groupDescription, &expirationSeconds)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	chat := r.buildGenericChat(chatID, chatType, createdAt, members, blockedBy,
		adminID, groupName, groupPhoto, groupDescription, expirationSeconds)

	return &chat, nil
}

func (r *GenericChatRepository) buildGenericChat(
	chatID uuid.UUID,
	chatType string,
	createdAt time.Time,
	members []uuid.UUID,
	blockedBy []uuid.UUID,
	adminID uuid.UUID,
	groupName string,
	groupPhoto string,
	groupDescription string,
	expirationSeconds *int,
) services.GenericChat {
	if chatType == services.ChatTypePersonal {
		return services.NewPersonalGenericChat(chatID, createdAt.Unix(), members, services.PersonalInfo{
			BlockedBy: blockedBy,
		})
	}

	if chatType == services.ChatTypeGroup {
		return services.NewGroupGenericChat(chatID, createdAt.Unix(), members, services.GroupInfo{
			Admin:            adminID,
			GroupName:        groupName,
			GroupDescription: groupDescription,
			GroupPhoto:       groupPhoto,
		})
	}

	if chatType == services.ChatTypeSecretPersonal {
		var exp *time.Duration
		if expirationSeconds != nil {
			cp := time.Duration(*expirationSeconds) * time.Second
			exp = &cp
		}
		return services.NewSecretPersonalGenericChat(chatID, createdAt.Unix(), members, services.SecretPersonalInfo{
			Expiration: exp,
		})
	}

	if chatType == services.ChatTypeSecretGroup {
		var exp *time.Duration
		if expirationSeconds != nil {
			cp := time.Duration(*expirationSeconds) * time.Second
			exp = &cp
		}
		return services.NewSecretGroupGenericChat(chatID, createdAt.Unix(), members, services.SecretGroupInfo{
			Admin:            adminID,
			GroupName:        groupName,
			GroupDescription: groupDescription,
			GroupPhoto:       groupPhoto,
			Expiration:       exp,
		})
	}

	panic(fmt.Errorf("unknown chat type is gotten from db: %s", chatType))
}

func (r *GenericChatRepository) GetChatType(ctx context.Context, id domain.ChatID) (string, error) {
	q := `SELECT chat_type FROM messaging.chat WHERE chat_id = $1`
	row := r.db.QueryRow(ctx, q, id)

	var chatType string
	if err := row.Scan(&chatType); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNotFound
		}
		return "", fmt.Errorf("getting chat type failed: %s", err)
	}

	return chatType, nil
}
