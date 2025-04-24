package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GenericChatRepository struct{}

func NewGenericChatRepository() *GenericChatRepository {
	return &GenericChatRepository{}
}

func (r *GenericChatRepository) GetByMemberID(
	ctx context.Context, db storage.ExecQuerier, memberID domain.UserID,
) ([]generic.Chat, error) {
	// I guess this query is so fucking slow in real production,
	// but for now we have more microservices than users, so nobody cares.
	//
	// TODO: optimize it.
	q := `
	SELECT 
		c.chat_id,
		c.chat_type,
		c.created_at,
		(SELECT ARRAY_AGG(me.user_id) 
		 FROM messaging.membership me
		 WHERE me.chat_id = c.chat_id),
		CASE WHEN c.chat_type = 'personal' 
				THEN (SELECT ARRAY_AGG(b.user_id) 
						  FROM messaging.blocking b 
						  WHERE b.chat_id = c.chat_id)
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
	WHERE m.user_id = $1`

	rows, err := db.Query(ctx, q, memberID)
	if err != nil {
		return nil, fmt.Errorf("getting chats by memberID failed: %s", err)
	}
	defer rows.Close()

	res := make([]generic.Chat, 0)
	for rows.Next() {
		var (
			chatID            uuid.UUID
			chatType          string
			createdAt         time.Time
			members           []uuid.UUID
			blockedBy         []uuid.UUID
			adminID           *uuid.UUID
			groupName         *string
			groupPhoto        *string
			groupDescription  *string
			expirationSeconds *int
		)
		err := rows.Scan(&chatID, &chatType, &createdAt, &members, &blockedBy,
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

func (r *GenericChatRepository) GetByChatID(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID,
) (*generic.Chat, error) {
	// I guess this query is so fucking slow in real production,
	// but for now we have more microservices than users, so nobody cares.
	//
	// TODO: optimize it.
	q := `
	SELECT 
		c.chat_id,
		c.chat_type,
		c.created_at,
		(SELECT ARRAY_AGG(me.user_id) 
		FROM messaging.membership me
		WHERE me.chat_id = c.chat_id),
		CASE WHEN c.chat_type = 'personal' 
				THEN (SELECT ARRAY_AGG(b.user_id) 
						FROM messaging.blocking b 
						WHERE b.chat_id = c.chat_id)
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
	WHERE c.chat_id = $1`

	row := db.QueryRow(ctx, q, id)

	var (
		chatID            uuid.UUID
		chatType          string
		createdAt         time.Time
		members           []uuid.UUID
		blockedBy         []uuid.UUID
		adminID           *uuid.UUID
		groupName         *string
		groupPhoto        *string
		groupDescription  *string
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

func deref[T any](t *T, defaultVal T) T {
	if t != nil {
		return *t
	}
	return defaultVal
}

func (r *GenericChatRepository) buildGenericChat(
	chatID uuid.UUID,
	chatType string,
	createdAt time.Time,
	members []uuid.UUID,
	blockedBy []uuid.UUID,
	adminID *uuid.UUID,
	groupName *string,
	groupPhoto *string,
	groupDescription *string,
	expirationSeconds *int,
) generic.Chat {
	result := generic.Chat{
		ChatID:    chatID,
		CreatedAt: createdAt.Unix(),
		Type:      chatType,
		Members:   members,
	}

	switch chatType {
	case domain.ChatTypePersonal:
		result.Info.Personal = &generic.PersonalInfo{
			BlockedBy: blockedBy,
		}
	case domain.ChatTypeGroup:
		result.Info.Group = &generic.GroupInfo{
			AdminID:     *adminID,
			Name:        *groupName,
			Description: deref(groupDescription, ""),
			GroupPhoto:  deref(groupPhoto, ""),
		}
	case domain.ChatTypeSecretPersonal:
		var exp *time.Duration
		if expirationSeconds != nil {
			cp := time.Duration(*expirationSeconds) * time.Second
			exp = &cp
		}
		result.Info.SecretPersonal = &generic.SecretPersonalInfo{
			Expiration: exp,
		}
	case domain.ChatTypeSecretGroup:
		var exp *time.Duration
		if expirationSeconds != nil {
			cp := time.Duration(*expirationSeconds) * time.Second
			exp = &cp
		}
		result.Info.SecretGroup = &generic.SecretGroupInfo{
			AdminID:     *adminID,
			Name:        *groupName,
			Description: deref(groupDescription, ""),
			GroupPhoto:  deref(groupPhoto, ""),
			Expiration:  exp,
		}
	default:
		panic(fmt.Errorf("unknown chat type is gotten from db: %s", chatType))
	}

	return result
}

func (r *GenericChatRepository) GetChatType(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID,
) (string, error) {
	q := `SELECT chat_type FROM messaging.chat WHERE chat_id = $1`
	row := db.QueryRow(ctx, q, id)

	var chatType string
	if err := row.Scan(&chatType); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNotFound
		}
		return "", fmt.Errorf("getting chat type failed: %s", err)
	}

	return chatType, nil
}
