package update

import (
	"context"
	"errors"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SecretUpdateRepository struct{}

func NewSecretUpdateRepository() *SecretUpdateRepository {
	return &SecretUpdateRepository{}
}

func (r *SecretUpdateRepository) CreateSecretUpdate(
	ctx context.Context, db storage.ExecQuerier, secret *domain.SecretUpdate,
) (*domain.SecretUpdate, error) {
	// Insert base update record
	q1 := `
	INSERT INTO messaging.update (chat_id, update_id, update_type, created_at, sender_id)
	VALUES ($1, $2, 'secret_update', $3, $4)
	RETURNING update_id`

	now := time.Now()
	var updateID int64
	err := db.QueryRow(ctx, q1,
		secret.ChatID,
		secret.UpdateID,
		now,
		uuid.UUID(secret.SenderID),
	).Scan(&updateID)
	if err != nil {
		return nil, err
	}
	secret.CreatedAt = domain.Timestamp(now.Unix())

	// Insert secret update specific data
	q2 := `
	INSERT INTO messaging.secret_update (chat_id, update_id, payload, key_hash, initialization_vector)
	VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(ctx, q2,
		secret.ChatID,
		updateID,
		secret.Data.Payload,
		secret.Data.KeyHash,
		secret.Data.Payload,
	)
	if err != nil {
		return nil, err
	}

	secret.UpdateID = domain.UpdateID(updateID)
	return secret, nil
}

func (r *SecretUpdateRepository) FindSecretUpdate(
	ctx context.Context, db storage.ExecQuerier, chatID domain.ChatID, updateID domain.UpdateID,
) (*domain.SecretUpdate, error) {
	q := `
	SELECT 
		u.created_at, 
		u.sender_id,
		s.payload,
		s.key_hash,
		s.initialization_vector
	FROM messaging.update u
	JOIN messaging.secret_update s ON u.chat_id = s.chat_id AND u.update_id = s.update_id
	WHERE u.chat_id = $1 AND u.update_id = $2`

	var (
		createdAt            time.Time
		senderID             uuid.UUID
		payload              []byte
		keyHash              []byte
		initializationVector []byte
	)

	err := db.QueryRow(ctx, q, chatID, updateID).Scan(
		&createdAt,
		&senderID,
		&payload,
		&keyHash,
		&initializationVector,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	// Get deletions for this update
	deletions, err := r.getDeletions(ctx, db, chatID, updateID)
	if err != nil {
		return nil, err
	}

	secretUpdate := &domain.SecretUpdate{
		Update: domain.Update{
			UpdateID:  updateID,
			ChatID:    chatID,
			SenderID:  domain.UserID(senderID),
			CreatedAt: domain.Timestamp(createdAt.Unix()),
			Deleted:   deletions,
		},
		Data: domain.SecretData{
			Payload: payload,
			KeyHash: domain.SecretKeyHash(keyHash),
			IV:      initializationVector,
		},
	}

	return secretUpdate, nil
}

func (r *SecretUpdateRepository) DeleteSecretUpdate(
	ctx context.Context, db storage.ExecQuerier, chatID domain.ChatID, updateID domain.UpdateID,
) error {
	q := `
	DELETE FROM messaging.update 
	WHERE chat_id = $1 AND update_id = $2 AND update_type = 'secret_update'`

	commandTag, err := db.Exec(ctx, q, chatID, updateID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *SecretUpdateRepository) CreateUpdateDeleted(
	ctx context.Context, db storage.ExecQuerier, deleted *domain.UpdateDeleted,
) (*domain.UpdateDeleted, error) {
	// Insert base update
	q1 := `
	INSERT INTO messaging.update (chat_id, update_id, update_type, created_at, sender_id)
	VALUES ($1, $2, 'update_deleted', $3, $4)
	RETURNING update_id`

	now := time.Now()
	var updateID int64
	err := db.QueryRow(ctx, q1,
		deleted.ChatID,
		deleted.UpdateID,
		now,
		uuid.UUID(deleted.SenderID),
	).Scan(&updateID)
	if err != nil {
		return nil, err
	}

	deleted.CreatedAt = domain.Timestamp(now.Unix())

	// Insert update_deleted specific data
	q2 := `
	INSERT INTO messaging.update_deleted_update (chat_id, update_id, reaction, deleted_update_id, mode)
	VALUES ($1, $2, '', $3, $4)`

	_, err = db.Exec(ctx, q2,
		deleted.ChatID,
		updateID,
		int64(deleted.DeletedID),
		deleted.Mode,
	)
	if err != nil {
		return nil, err
	}

	deleted.UpdateID = domain.UpdateID(updateID)
	return deleted, nil
}

// Helper function to get deletions for an update
func (r *SecretUpdateRepository) getDeletions(
	ctx context.Context, db storage.ExecQuerier, chatID domain.ChatID, updateID domain.UpdateID,
) ([]*domain.UpdateDeleted, error) {
	q := `
	SELECT 
		dud.update_id, u.created_at, u.sender_id,
		dud.deleted_update_id, dud.mode
	FROM messaging.update_deleted_update dud
	JOIN messaging.update u ON dud.chat_id = u.chat_id AND dud.update_id = u.update_id
	WHERE dud.deleted_update_id = $1 AND dud.chat_id = $2`

	rows, err := db.Query(ctx, q, updateID, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deletions []*domain.UpdateDeleted
	for rows.Next() {
		var (
			id              int64
			createdAt       time.Time
			senderID        uuid.UUID
			deletedUpdateID int64
			modeStr         string
		)

		if err := rows.Scan(&id, &createdAt, &senderID, &deletedUpdateID, &modeStr); err != nil {
			return nil, err
		}

		deletion := &domain.UpdateDeleted{
			Update: domain.Update{
				UpdateID:  domain.UpdateID(id),
				ChatID:    chatID,
				SenderID:  domain.UserID(senderID),
				CreatedAt: domain.Timestamp(createdAt.Unix()),
				Deleted:   nil, // We don't fetch nested deletions
			},
			DeletedID: domain.UpdateID(deletedUpdateID),
			Mode:      domain.DeleteMode(modeStr),
		}

		deletions = append(deletions, deletion)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deletions, nil
}
