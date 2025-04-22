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

type UpdateRepository struct{}

func NewUpdateRepository() *UpdateRepository {
	return &UpdateRepository{}
}

func (r *UpdateRepository) FindGenericMessage(
	ctx context.Context, db storage.ExecQuerier, chatID domain.ChatID, updateID domain.UpdateID,
) (*domain.Message, error) {
	q := `
	SELECT 
		u.update_type,
		u.created_at,
		u.sender_id,
		COALESCE(tm.reply_to_id, fm.reply_to_id)
	FROM messaging.update u
		LEFT JOIN messaging.text_message_update tm ON tm.chat_id = u.chat_id AND tm.update_id = u.update_id
		LEFT JOIN messaging.file_message_update fm ON fm.chat_id = u.chat_id AND fm.update_id = u.update_id
	WHERE u.chat_id = $1 AND u.update_id = $2`

	var (
		updateType string
		createdAt  time.Time
		senderID   uuid.UUID
		replyToID  *int64
	)

	err := db.QueryRow(ctx, q, chatID, updateID).Scan(
		&updateType,
		&createdAt,
		&senderID,
		&replyToID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	deletions, err := r.getDeletions(ctx, db, chatID, updateID)
	if err != nil {
		return nil, err
	}

	var replyToUpdateID *domain.UpdateID
	if replyToID != nil {
		id := domain.UpdateID(*replyToID)
		replyToUpdateID = &id
	}

	message := &domain.Message{
		Update: domain.Update{
			UpdateID:  updateID,
			ChatID:    chatID,
			SenderID:  domain.UserID(senderID),
			CreatedAt: domain.Timestamp(createdAt.Unix()),
			Deleted:   deletions,
		},
		ReplyTo:   replyToUpdateID,
		Forwarded: false, // This isn't in the schema, would need additional data
	}

	return message, nil
}

func (r *UpdateRepository) DeleteUpdate(
	ctx context.Context, db storage.ExecQuerier, chatID domain.ChatID, updateID domain.UpdateID,
) error {
	q := `
	DELETE FROM messaging.update 
	WHERE chat_id = $1 AND update_id = $2`

	_, err := db.Exec(ctx, q, chatID, updateID)
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrNotFound
	}
	return err
}

func (r *UpdateRepository) CreateUpdateDeleted(
	ctx context.Context, db storage.ExecQuerier, deleted *domain.UpdateDeleted,
) (*domain.UpdateDeleted, error) {
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

	q2 := `
	INSERT INTO messaging.update_deleted_update (chat_id, update_id, deleted_update_id, mode)
	VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(ctx, q2,
		deleted.ChatID,
		updateID,
		deleted.DeletedID,
		deleted.Mode,
	)
	if err != nil {
		return nil, err
	}

	deleted.UpdateID = domain.UpdateID(updateID)
	return deleted, nil
}

func (r *UpdateRepository) CreateTextMessage(
	ctx context.Context, db storage.ExecQuerier, msg *domain.TextMessage,
) (*domain.TextMessage, error) {
	q1 := `
	INSERT INTO messaging.update (chat_id, update_id, update_type, created_at, sender_id)
	VALUES ($1, $2, 'text_message', $3, $4)
	RETURNING update_id`

	now := time.Now()
	var updateID int64
	err := db.QueryRow(ctx, q1,
		msg.ChatID,
		msg.UpdateID,
		now,
		uuid.UUID(msg.SenderID),
	).Scan(&updateID)
	if err != nil {
		return nil, err
	}
	msg.CreatedAt = domain.Timestamp(now.Unix())

	// Insert text_message specific data
	q2 := `
	INSERT INTO messaging.text_message_update (chat_id, update_id, text, reply_to_id)
	VALUES ($1, $2, $3, $4)`

	var replyToID *int64
	if msg.ReplyTo != nil {
		id := int64(*msg.ReplyTo)
		replyToID = &id
	}

	_, err = db.Exec(ctx, q2,
		msg.ChatID,
		updateID,
		msg.Text,
		replyToID,
	)
	if err != nil {
		return nil, err
	}

	msg.UpdateID = domain.UpdateID(updateID)
	return msg, nil
}

func (r *UpdateRepository) CreateTextMessageEdited(
	ctx context.Context, db storage.ExecQuerier, edited *domain.TextMessageEdited,
) (*domain.TextMessageEdited, error) {
	// Insert base update
	q1 := `
	INSERT INTO messaging.update (chat_id, update_id, update_type, created_at, sender_id)
	VALUES ($1, $2, 'text_message_edited', $3, $4)
	RETURNING update_id`

	now := time.Now()
	var updateID int64
	err := db.QueryRow(ctx, q1,
		edited.ChatID,
		edited.UpdateID,
		now,
		uuid.UUID(edited.SenderID),
	).Scan(&updateID)
	if err != nil {
		return nil, err
	}
	edited.CreatedAt = domain.Timestamp(now.Unix())

	// Insert text_message_edited specific data
	q2 := `
	INSERT INTO messaging.text_message_edited_update (chat_id, update_id, new_text, message_id)
	VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(ctx, q2,
		edited.ChatID,
		updateID,
		edited.NewText,
		int64(edited.MessageID),
	)
	if err != nil {
		return nil, err
	}

	edited.UpdateID = domain.UpdateID(updateID)
	return edited, nil
}

func (r *UpdateRepository) FindTextMessage(
	ctx context.Context, db storage.ExecQuerier, chatID domain.ChatID, updateID domain.UpdateID,
) (*domain.TextMessage, error) {
	// Get base message first
	message, err := r.FindGenericMessage(ctx, db, chatID, updateID)
	if err != nil {
		return nil, err
	}

	// Get text message specific data
	q := `
	SELECT 
		tm.text 
	FROM messaging.text_message_update tm
	WHERE tm.chat_id = $1 AND tm.update_id = $2`

	var text string
	err = db.QueryRow(ctx, q, chatID, updateID).Scan(&text)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	// Check if there's an edit for this message
	edits, err := r.getTextMessageEdits(ctx, db, chatID, updateID)
	if err != nil {
		return nil, err
	}

	var edited *domain.TextMessageEdited
	if len(edits) > 0 {
		// Use the most recent edit
		edited = edits[0]
	}

	textMsg := &domain.TextMessage{
		Message: *message,
		Text:    text,
		Edited:  edited,
	}

	return textMsg, nil
}

func (r *UpdateRepository) UpdateTextMessage(
	ctx context.Context, db storage.ExecQuerier, msg *domain.TextMessage,
) (*domain.TextMessage, error) {
	q := `
	UPDATE messaging.text_message_update
	SET text = $3, reply_to_id = $4
	WHERE chat_id = $1 AND update_id = $2`

	var replyToID *int64
	if msg.ReplyTo != nil {
		id := int64(*msg.ReplyTo)
		replyToID = &id
	}

	_, err := db.Exec(ctx, q,
		msg.ChatID,
		msg.UpdateID,
		msg.Text,
		replyToID,
	)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *UpdateRepository) CreateReaction(
	ctx context.Context, db storage.ExecQuerier, reaction *domain.Reaction,
) (*domain.Reaction, error) {
	// Insert base update
	q1 := `
	INSERT INTO messaging.update (chat_id, update_id, update_type, created_at, sender_id)
	VALUES ($1, $2, 'reaction', $3, $4)
	RETURNING update_id`

	now := time.Now()
	var updateID int64
	err := db.QueryRow(ctx, q1,
		reaction.ChatID,
		reaction.UpdateID,
		now,
		uuid.UUID(reaction.SenderID),
	).Scan(&updateID)
	if err != nil {
		return nil, err
	}

	// Insert reaction specific data
	q2 := `
	INSERT INTO messaging.reaction_update (chat_id, update_id, reaction, message_id)
	VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(ctx, q2,
		reaction.ChatID,
		updateID,
		reaction.Type,
		reaction.MessageID,
	)
	if err != nil {
		return nil, err
	}

	reaction.UpdateID = domain.UpdateID(updateID)
	return reaction, nil
}

func (r *UpdateRepository) FindReaction(
	ctx context.Context, db storage.ExecQuerier, chatID domain.ChatID, updateID domain.UpdateID,
) (*domain.Reaction, error) {
	q := `
	SELECT 
		u.created_at, 
		u.sender_id,
		r.reaction,
		r.message_id
	FROM messaging.update u
	JOIN messaging.reaction_update r ON u.chat_id = r.chat_id AND u.update_id = r.update_id
	WHERE u.chat_id = $1 AND u.update_id = $2`

	var (
		createdAt    time.Time
		senderID     uuid.UUID
		reactionType string
		messageID    int64
	)

	err := db.QueryRow(ctx, q, chatID, updateID).Scan(
		&createdAt,
		&senderID,
		&reactionType,
		&messageID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	deletions, err := r.getDeletions(ctx, db, chatID, updateID)
	if err != nil {
		return nil, err
	}

	reaction := &domain.Reaction{
		Update: domain.Update{
			UpdateID:  updateID,
			ChatID:    chatID,
			SenderID:  domain.UserID(senderID),
			CreatedAt: domain.Timestamp(createdAt.Unix()),
			Deleted:   deletions,
		},
		Type:      domain.ReactionType(reactionType),
		MessageID: domain.UpdateID(messageID),
	}

	return reaction, nil
}

func (r *UpdateRepository) FindFileMessage(
	ctx context.Context, db storage.ExecQuerier, chatID domain.ChatID, updateID domain.UpdateID,
) (*domain.FileMessage, error) {
	// Get base message first
	message, err := r.FindGenericMessage(ctx, db, chatID, updateID)
	if err != nil {
		return nil, err
	}

	// Get file message specific data
	q := `
	SELECT 
		file_id, file_name, file_mime_type, file_size, file_url, file_created_at
	FROM messaging.file_message_update 
	WHERE chat_id = $1 AND update_id = $2`

	var (
		fileID        uuid.UUID
		fileName      string
		fileMimeType  string
		fileSize      int64
		fileURL       string
		fileCreatedAt int64
	)

	err = db.QueryRow(ctx, q, chatID, updateID).Scan(
		&fileID,
		&fileName,
		&fileMimeType,
		&fileSize,
		&fileURL,
		&fileCreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	fileMsg := &domain.FileMessage{
		Message: *message,
		File: domain.FileMeta{
			FileId:    fileID,
			FileName:  fileName,
			MimeType:  fileMimeType,
			FileSize:  fileSize,
			FileURL:   domain.URL(fileURL),
			CreatedAt: domain.Timestamp(fileCreatedAt),
		},
	}

	return fileMsg, nil
}

func (r *UpdateRepository) CreateFileMessage(
	ctx context.Context, db storage.ExecQuerier, msg *domain.FileMessage,
) (*domain.FileMessage, error) {
	// Insert base update
	q1 := `
	INSERT INTO messaging.update (chat_id, update_id, update_type, created_at, sender_id)
	VALUES ($1, $2, 'file_message', $3, $4)
	RETURNING update_id`

	now := time.Now()
	var updateID int64
	err := db.QueryRow(ctx, q1,
		msg.ChatID,
		msg.UpdateID,
		now,
		uuid.UUID(msg.SenderID),
	).Scan(&updateID)
	if err != nil {
		return nil, err
	}
	msg.CreatedAt = domain.Timestamp(now.Unix())

	// Insert file_message specific data
	q2 := `
	INSERT INTO messaging.file_message_update (
		chat_id, update_id, file_id, file_name, file_mime_type, 
		file_size, file_url, file_created_at, reply_to_id
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	var replyToID *int64
	if msg.ReplyTo != nil {
		id := int64(*msg.ReplyTo)
		replyToID = &id
	}

	_, err = db.Exec(ctx, q2,
		msg.ChatID,
		updateID,
		msg.File.FileId,
		msg.File.FileName,
		msg.File.MimeType,
		msg.File.FileSize,
		string(msg.File.FileURL),
		msg.File.CreatedAt,
		replyToID,
	)
	if err != nil {
		return nil, err
	}

	msg.UpdateID = domain.UpdateID(updateID)
	return msg, nil
}

// Helper functions

func (r *UpdateRepository) getDeletions(
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

func (r *UpdateRepository) getTextMessageEdits(
	ctx context.Context, db storage.ExecQuerier, chatID domain.ChatID, messageID domain.UpdateID,
) ([]*domain.TextMessageEdited, error) {
	q := `
	SELECT 
		e.update_id, u.created_at, u.sender_id, e.new_text
	FROM messaging.text_message_edited_update e
	JOIN messaging.update u ON e.chat_id = u.chat_id AND e.update_id = u.update_id
	WHERE e.message_id = $1 AND e.chat_id = $2
	ORDER BY u.created_at DESC`

	rows, err := db.Query(ctx, q, messageID, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edits []*domain.TextMessageEdited
	for rows.Next() {
		var (
			id        int64
			createdAt time.Time
			senderID  uuid.UUID
			newText   string
		)

		if err := rows.Scan(&id, &createdAt, &senderID, &newText); err != nil {
			return nil, err
		}

		deletions, err := r.getDeletions(ctx, db, chatID, domain.UpdateID(id))
		if err != nil {
			return nil, err
		}

		edit := &domain.TextMessageEdited{
			Update: domain.Update{
				UpdateID:  domain.UpdateID(id),
				ChatID:    chatID,
				SenderID:  domain.UserID(senderID),
				CreatedAt: domain.Timestamp(createdAt.Unix()),
				Deleted:   deletions,
			},
			MessageID: messageID,
			NewText:   newText,
		}

		edits = append(edits, edit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return edits, nil
}
