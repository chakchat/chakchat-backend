package update

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type GenericUpdateRepository struct{}

func NewGenericUpdateRepository() *GenericUpdateRepository {
	return &GenericUpdateRepository{}
}

func (r *GenericUpdateRepository) GetRange(
	ctx context.Context,
	db storage.ExecQuerier,
	visibleTo domain.UserID,
	chatID domain.ChatID,
	from, to domain.UpdateID,
) ([]services.GenericUpdate, error) {
	// TODO: This query should be optimized.
	// Especially the check of being deleted.
	q := `
	SELECT DISTINCT
		u.update_id,
		u.update_type,
		u.created_at,
		u.sender_id
	FROM messaging.update u
		LEFT JOIN messaging.update_deleted_update ud ON ud.deleted_update_id = u.update_id AND ud.chat_id = u.chat_id
	WHERE u.chat_id = $1 
		AND u.update_id BETWEEN $3 AND $4
		AND ud.mode IS DISTINCT FROM 'for_all'
		AND (
			ud.mode IS DISTINCT FROM 'for_deletion_sender' 
			OR $2 <> (
				SELECT sender_id 
				FROM messaging.update 
				WHERE chat_id = $1 
					AND update_id = ud.update_id)
		)`

	rows, err := db.Query(ctx, q, chatID, visibleTo, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updates := make([]services.GenericUpdate, 0, to-from+1)

	for rows.Next() {
		var (
			updateID   int64
			updateType string
			createdAt  time.Time
			senderID   uuid.UUID
		)
		if err := rows.Scan(&updateID, &updateType, &createdAt, &senderID); err != nil {
			return nil, err
		}
		updates = append(updates, services.GenericUpdate{
			UpdateID:   updateID,
			ChatID:     uuid.UUID(chatID),
			SenderID:   senderID,
			UpdateType: updateType,
			CreatedAt:  createdAt.Unix(),
			Info:       services.GenericUpdateInfo{},
		})
	}

	// Fill in details for each update type
	if err := r.fillTextMessages(ctx, db, chatID, updates); err != nil {
		return nil, err
	}

	if err := r.fillTextMessageEdited(ctx, db, chatID, updates); err != nil {
		return nil, err
	}

	if err := r.fillFileMessages(ctx, db, chatID, updates); err != nil {
		return nil, err
	}

	if err := r.fillReactions(ctx, db, chatID, updates); err != nil {
		return nil, err
	}

	if err := r.fillUpdateDeleted(ctx, db, chatID, updates); err != nil {
		return nil, err
	}

	if err := r.fillSecretUpdates(ctx, db, chatID, updates); err != nil {
		return nil, err
	}

	return updates, nil
}

func (r *GenericUpdateRepository) Get(
	ctx context.Context,
	db storage.ExecQuerier,
	visibleTo domain.UserID,
	chatID domain.ChatID,
	updateID domain.UpdateID,
) (*services.GenericUpdate, error) {
	q := `
	SELECT
		u.update_id,
		u.update_type,
		u.created_at,
		u.sender_id
	FROM messaging.update u
		LEFT JOIN messaging.update_deleted_update ud ON ud.deleted_update_id = u.update_id AND ud.chat_id = u.chat_id
	WHERE u.chat_id = $1 
		AND u.update_id = $2
		AND ud.mode IS DISTINCT FROM 'for_all'
		AND (
			ud.mode IS DISTINCT FROM 'for_deletion_sender' 
			OR $3 <> (
				SELECT sender_id 
				FROM messaging.update 
				WHERE chat_id = $1 
					AND update_id = ud.update_id)
		)`

	var (
		updateIDVal int64
		updateType  string
		createdAt   time.Time
		senderID    uuid.UUID
	)

	err := db.QueryRow(ctx, q, chatID, updateID, visibleTo).Scan(
		&updateIDVal,
		&updateType,
		&createdAt,
		&senderID,
	)
	if err != nil {
		return nil, err
	}

	update := services.GenericUpdate{
		UpdateID:   updateIDVal,
		ChatID:     uuid.UUID(chatID),
		SenderID:   senderID,
		UpdateType: updateType,
		CreatedAt:  createdAt.Unix(),
		Info:       services.GenericUpdateInfo{},
	}

	// Fill in details based on update type
	updates := []services.GenericUpdate{update}

	switch updateType {
	case services.UpdateTypeTextMessage:
		if err := r.fillTextMessages(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case services.UpdateTypeTextMessageEdited:
		if err := r.fillTextMessageEdited(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case services.UpdateTypeFileMessage:
		if err := r.fillFileMessages(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case services.UpdateTypeReaction:
		if err := r.fillReactions(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case services.UpdateTypeDeleted:
		if err := r.fillUpdateDeleted(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case services.UpdateTypeSecret:
		if err := r.fillSecretUpdates(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	}

	return &updates[0], nil
}

func (r *GenericUpdateRepository) fillTextMessages(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []services.GenericUpdate,
) error {
	ids := updateTypesIDs(updates, services.UpdateTypeTextMessage)
	if len(ids) == 0 {
		return nil
	}

	q := fmt.Sprintf(`
	SELECT
		tm.update_id,
		tm.text,
		tm.reply_to_id
	FROM messaging.text_message_update tm
	WHERE tm.chat_id = $1 
		AND tm.update_id IN %s
	`, sqlArgsArr(2, len(ids)))

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...))
	if err != nil {
		return err
	}
	defer rows.Close()

	// Map to store text messages by update_id
	textMsgs := make(map[int64]services.TextMessageInfo)

	for rows.Next() {
		var (
			updateID  int64
			text      string
			replyToID *int64
		)

		if err := rows.Scan(&updateID, &text, &replyToID); err != nil {
			return err
		}

		textMsgs[updateID] = services.TextMessageInfo{
			Text:    text,
			ReplyTo: replyToID,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Get edited info for text messages
	editedInfo, err := r.getTextEdits(ctx, db, chatID, ids)
	if err != nil {
		return err
	}

	// Update the updates with the text message info
	for i, update := range updates {
		if update.UpdateType == services.UpdateTypeTextMessage {
			if info, ok := textMsgs[update.UpdateID]; ok {
				// Check for edits
				if edit, exists := editedInfo[update.UpdateID]; exists {
					info.Edited = edit
				}

				updates[i].Info.TextMessage = &info
			}
		}
	}

	return nil
}

// Helper function to get edits for text messages
func (r *GenericUpdateRepository) getTextEdits(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	messageIDs []int64,
) (map[int64]*services.GenericUpdate, error) {
	if len(messageIDs) == 0 {
		return nil, nil
	}

	q := fmt.Sprintf(`
	SELECT
		tme.update_id,
		tme.message_id,
		tme.new_text,
		u.created_at,
		u.sender_id
	FROM messaging.text_message_edited_update tme
	JOIN messaging.update u ON u.chat_id = tme.chat_id AND u.update_id = tme.update_id
	WHERE tme.chat_id = $1 
		AND tme.message_id IN %s
	ORDER BY u.created_at DESC
	`, sqlArgsArr(2, len(messageIDs)))

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(messageIDs)...))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map message ID to the most recent edit
	result := make(map[int64]*services.GenericUpdate)

	for rows.Next() {
		var (
			updateID  int64
			messageID int64
			newText   string
			createdAt time.Time
			senderID  uuid.UUID
		)

		if err := rows.Scan(&updateID, &messageID, &newText, &createdAt, &senderID); err != nil {
			return nil, err
		}

		// Only store the first edit we encounter for each message ID (most recent)
		if _, exists := result[messageID]; !exists {
			editInfo := services.TextMessageEditedInfo{
				MessageID: messageID,
				NewText:   newText,
			}

			result[messageID] = &services.GenericUpdate{
				UpdateID:   updateID,
				ChatID:     uuid.UUID(chatID),
				SenderID:   senderID,
				UpdateType: services.UpdateTypeTextMessageEdited,
				CreatedAt:  createdAt.Unix(),
				Info: services.GenericUpdateInfo{
					TextMessageEdited: &editInfo,
				},
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *GenericUpdateRepository) fillTextMessageEdited(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []services.GenericUpdate,
) error {
	ids := updateTypesIDs(updates, services.UpdateTypeTextMessageEdited)
	if len(ids) == 0 {
		return nil
	}

	q := fmt.Sprintf(`
	SELECT
		tme.update_id,
		tme.new_text,
		tme.message_id
	FROM messaging.text_message_edited_update tme
	WHERE tme.chat_id = $1 
		AND tme.update_id IN %s
	`, sqlArgsArr(2, len(ids)))

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...))
	if err != nil {
		return err
	}
	defer rows.Close()

	edits := make(map[int64]services.TextMessageEditedInfo)

	for rows.Next() {
		var (
			updateID  int64
			newText   string
			messageID int64
		)

		if err := rows.Scan(&updateID, &newText, &messageID); err != nil {
			return err
		}

		edits[updateID] = services.TextMessageEditedInfo{
			NewText:   newText,
			MessageID: messageID,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the edit info
	for i, update := range updates {
		if update.UpdateType == services.UpdateTypeTextMessageEdited {
			if info, ok := edits[update.UpdateID]; ok {
				updates[i].Info.TextMessageEdited = &info
			}
		}
	}

	return nil
}

func (r *GenericUpdateRepository) fillFileMessages(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []services.GenericUpdate,
) error {
	ids := updateTypesIDs(updates, services.UpdateTypeFileMessage)
	if len(ids) == 0 {
		return nil
	}

	q := fmt.Sprintf(`
	SELECT
		fm.update_id,
		fm.file_id,
		fm.file_name,
		fm.file_mime_type,
		fm.file_size,
		fm.file_url,
		fm.file_created_at,
		fm.reply_to_id
	FROM messaging.file_message_update fm
	WHERE fm.chat_id = $1 
		AND fm.update_id IN %s
	`, sqlArgsArr(2, len(ids)))

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...))
	if err != nil {
		return err
	}
	defer rows.Close()

	fileMsgs := make(map[int64]services.FileMessageInfo)

	for rows.Next() {
		var (
			updateID      int64
			fileID        uuid.UUID
			fileName      string
			fileMimeType  string
			fileSize      int64
			fileURL       string
			fileCreatedAt time.Time
			replyToID     *int64
		)

		if err := rows.Scan(
			&updateID,
			&fileID,
			&fileName,
			&fileMimeType,
			&fileSize,
			&fileURL,
			&fileCreatedAt,
			&replyToID,
		); err != nil {
			return err
		}

		fileMsgs[updateID] = services.FileMessageInfo{
			File: dto.FileMetaDTO{
				FileId:    fileID,
				FileName:  fileName,
				MimeType:  fileMimeType,
				FileSize:  fileSize,
				FileURL:   fileURL,
				CreatedAt: fileCreatedAt.Unix(),
			},
			ReplyTo: replyToID,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the file message info
	for i, update := range updates {
		if update.UpdateType == services.UpdateTypeFileMessage {
			if info, ok := fileMsgs[update.UpdateID]; ok {
				updates[i].Info.FileMessage = &info
			}
		}
	}

	return nil
}

func (r *GenericUpdateRepository) fillReactions(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []services.GenericUpdate,
) error {
	ids := updateTypesIDs(updates, services.UpdateTypeReaction)
	if len(ids) == 0 {
		return nil
	}

	q := fmt.Sprintf(`
	SELECT
		r.update_id,
		r.reaction,
		r.message_id
	FROM messaging.reaction_update r
	WHERE r.chat_id = $1 
		AND r.update_id IN %s
	`, sqlArgsArr(2, len(ids)))

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...))
	if err != nil {
		return err
	}
	defer rows.Close()

	reactions := make(map[int64]services.ReactionInfo)

	for rows.Next() {
		var (
			updateID  int64
			reaction  string
			messageID int64
		)

		if err := rows.Scan(&updateID, &reaction, &messageID); err != nil {
			return err
		}

		reactions[updateID] = services.ReactionInfo{
			Reaction: reaction,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the reaction info
	for i, update := range updates {
		if update.UpdateType == services.UpdateTypeReaction {
			if info, ok := reactions[update.UpdateID]; ok {
				updates[i].Info.Reaction = &info
			}
		}
	}

	return nil
}

func (r *GenericUpdateRepository) fillUpdateDeleted(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []services.GenericUpdate,
) error {
	ids := updateTypesIDs(updates, services.UpdateTypeDeleted)
	if len(ids) == 0 {
		return nil
	}

	q := fmt.Sprintf(`
	SELECT
		ud.update_id,
		ud.deleted_update_id,
		ud.mode
	FROM messaging.update_deleted_update ud
	WHERE ud.chat_id = $1 
		AND ud.update_id IN %s
	`, sqlArgsArr(2, len(ids)))

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...))
	if err != nil {
		return err
	}
	defer rows.Close()

	deletedUpdates := make(map[int64]services.DeletedInfo)

	for rows.Next() {
		var (
			updateID        int64
			deletedUpdateID int64
			mode            string
		)

		if err := rows.Scan(&updateID, &deletedUpdateID, &mode); err != nil {
			return err
		}

		deletedUpdates[updateID] = services.DeletedInfo{
			DeletedID:  deletedUpdateID,
			DeleteMode: mode,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the deleted update info
	for i, update := range updates {
		if update.UpdateType == services.UpdateTypeDeleted {
			if info, ok := deletedUpdates[update.UpdateID]; ok {
				updates[i].Info.Deleted = &info
			}
		}
	}

	return nil
}

func (r *GenericUpdateRepository) fillSecretUpdates(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []services.GenericUpdate,
) error {
	ids := updateTypesIDs(updates, services.UpdateTypeSecret)
	if len(ids) == 0 {
		return nil
	}

	q := fmt.Sprintf(`
	SELECT
		su.update_id,
		su.payload,
		su.key_hash,
		su.initialization_vector
	FROM messaging.secret_update su
	WHERE su.chat_id = $1 
		AND su.update_id IN %s
	`, sqlArgsArr(2, len(ids)))

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...))
	if err != nil {
		return err
	}
	defer rows.Close()

	secretUpdates := make(map[int64]services.SecretUpdateInfo)

	for rows.Next() {
		var (
			updateID             int64
			payload              []byte
			keyHash              []byte
			initializationVector []byte
		)

		if err := rows.Scan(&updateID, &payload, &keyHash, &initializationVector); err != nil {
			return err
		}

		secretUpdates[updateID] = services.SecretUpdateInfo{
			Payload:              payload,
			KeyHash:              keyHash,
			InitializationVector: initializationVector,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the secret update info
	for i, update := range updates {
		if update.UpdateType == services.UpdateTypeSecret {
			if info, ok := secretUpdates[update.UpdateID]; ok {
				updates[i].Info.Secret = &info
			}
		}
	}

	return nil
}

// Helper function to convert []int64 to []any for SQL parameters
func idsToAny(ids []int64) []any {
	result := make([]any, len(ids))
	for i, id := range ids {
		result[i] = id
	}
	return result
}

func updateTypesIDs(updates []services.GenericUpdate, typ string) []int64 {
	ids := make([]int64, 0)
	for _, up := range updates {
		if up.UpdateType == typ {
			ids = append(ids, up.UpdateID)
		}
	}
	return ids
}

func sqlArgsArr(start, cnt int) string {
	var s strings.Builder
	s.WriteRune('(')
	for i := start; i < start+cnt; i++ {
		s.WriteRune('$')
		s.WriteString(strconv.Itoa(i))
		if i != start+cnt-1 {
			s.WriteString(", ")
		}
	}
	s.WriteRune(')')
	return s.String()
}
