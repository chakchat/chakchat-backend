package update

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GenericUpdateRepository struct{}

func NewGenericUpdateRepository() *GenericUpdateRepository {
	return &GenericUpdateRepository{}
}

func (r *GenericUpdateRepository) GetLastUpdateID(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID,
) (domain.UpdateID, error) {
	q := `
	SELECT last_update_id
	FROM messaging.chat_sequence
	WHERE chat_id = $1`

	var lastUpdateID int64
	if err := db.QueryRow(ctx, q, id).Scan(&lastUpdateID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return domain.UpdateID(lastUpdateID), nil
}

func (r *GenericUpdateRepository) GetRange(
	ctx context.Context,
	db storage.ExecQuerier,
	visibleTo domain.UserID,
	chatID domain.ChatID,
	from, to domain.UpdateID,
) ([]generic.Update, error) {
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

	updates := make([]generic.Update, 0, max(0, to-from+1))

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
		updates = append(updates, generic.Update{
			UpdateID:   updateID,
			ChatID:     uuid.UUID(chatID),
			SenderID:   senderID,
			UpdateType: updateType,
			CreatedAt:  createdAt.Unix(),
			Content:    generic.UpdateContent{},
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
) (*generic.Update, error) {
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

	update := generic.Update{
		UpdateID:   updateIDVal,
		ChatID:     uuid.UUID(chatID),
		SenderID:   senderID,
		UpdateType: updateType,
		CreatedAt:  createdAt.Unix(),
		Content:    generic.UpdateContent{},
	}

	// Fill in details based on update type
	updates := []generic.Update{update}

	switch updateType {
	case domain.UpdateTypeTextMessage:
		if err := r.fillTextMessages(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case domain.UpdateTypeTextMessageEdited:
		if err := r.fillTextMessageEdited(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case domain.UpdateTypeFileMessage:
		if err := r.fillFileMessages(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case domain.UpdateTypeReaction:
		if err := r.fillReactions(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case domain.UpdateTypeDeleted:
		if err := r.fillUpdateDeleted(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	case domain.UpdateTypeSecret:
		if err := r.fillSecretUpdates(ctx, db, chatID, updates); err != nil {
			return nil, err
		}
	}

	return &updates[0], nil
}

func (r *GenericUpdateRepository) FetchLast( // TODO: rename to latest
	ctx context.Context,
	db storage.ExecQuerier,
	visibleTo domain.UserID,
	chatID domain.ChatID,
	opts ...repository.FetchLastOption,
) ([]generic.Update, error) {
	opt := repository.NewFetchLastOptions(opts...)

	var updateTypes []string
	if opt.Mode & repository.FetchLastModeMessages != 0 {
		updateTypes = append(updateTypes, "text_message", "file_message")
	}
	if opt.Mode & repository.FetchLastModeReactions != 0 {
		updateTypes = append(updateTypes, "reaction")
	}
	if opt.Mode & repository.FetchLastModeSecret != 0 {
		updateTypes = append(updateTypes, "secret_update")
	}

	lo, found, err := r.getUpdateIDFromLast(
		ctx, db, uuid.UUID(visibleTo), uuid.UUID(chatID), updateTypes, opt.Count,
	)
	if err != nil {
		return nil, err
	}
	if !found {
		return make([]generic.Update, 0), nil
	}

	hi, err := r.GetLastUpdateID(ctx, db, chatID)
	if err != nil {
		return nil, err
	}

	return r.GetRange(ctx, db, visibleTo, chatID, lo, hi)
}

func (r *GenericUpdateRepository) getUpdateIDFromLast(
	ctx context.Context,
	db storage.ExecQuerier,
	visibleTo, chatID uuid.UUID,
	updateTypes []string,
	count int,
) (_ domain.UpdateID, found bool, _ error) {
	var additionalWhereClause string
	if updateTypes != nil {
		for i := range updateTypes {
			updateTypes[i] = fmt.Sprintf(`'%s'`, updateTypes[i])
		}
		additionalWhereClause = fmt.Sprintf("AND u.update_type IN (%s)", strings.Join(updateTypes, ", "))
	}
	// TODO: Please use squirrel package if such cringe will be repeated
	q := fmt.Sprintf(`
	SELECT u.update_id
	FROM messaging.update u
		LEFT JOIN messaging.update_deleted_update ud 
			ON ud.deleted_update_id = u.update_id AND ud.chat_id = u.chat_id
	WHERE u.chat_id = $1
		%s -- Here is check for update type.
		AND ud.mode IS DISTINCT FROM 'for_all'
		AND (
			ud.mode IS DISTINCT FROM 'for_deletion_sender' 
			OR $2 <> (
				SELECT u.sender_id 
				FROM messaging.update 
				WHERE u.chat_id = $1 
					AND u.update_id = ud.update_id)
		)
	ORDER BY u.update_id DESC
	LIMIT $3
	`, additionalWhereClause)

	rows, err := db.Query(ctx, q, chatID, visibleTo, count)
	if err != nil {
		return 0, false, err
	}
	defer rows.Close()

	var minUpdateID *int64
	for rows.Next() {
		var updateID int64
		if err := rows.Scan(&updateID); err != nil {
			return 0, false, err
		}
		if minUpdateID == nil || *minUpdateID > updateID {
			minUpdateID = &updateID
		}
	}
	if err := rows.Err(); err != nil {
		return 0, false, err
	}
	if minUpdateID == nil {
		return 0, false, nil
	}

	return domain.UpdateID(*minUpdateID), true, nil
}

func (r *GenericUpdateRepository) fillTextMessages(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []generic.Update,
) error {
	ids := updateTypesIDs(updates, domain.UpdateTypeTextMessage)
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

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...)...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Map to store text messages by update_id
	textMsgs := make(map[int64]generic.TextMessageContent)

	for rows.Next() {
		var (
			updateID  int64
			text      string
			replyToID *int64
		)

		if err := rows.Scan(&updateID, &text, &replyToID); err != nil {
			return err
		}

		textMsgs[updateID] = generic.TextMessageContent{
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
		if update.UpdateType == domain.UpdateTypeTextMessage {
			if info, ok := textMsgs[update.UpdateID]; ok {
				// Check for edits
				if edit, exists := editedInfo[update.UpdateID]; exists {
					info.Edited = edit
				}

				reactions, err := r.getMessageReactions(ctx, db, chatID, domain.UpdateID(update.UpdateID))
				if err != nil {
					return err
				}
				if reactions != nil {
					info.Reactions = reactions
				}

				updates[i].Content.TextMessage = &info
			} else {
				panic("database is inconsistent!!!")
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
) (map[int64]*generic.Update, error) {
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

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(messageIDs)...)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map message ID to the most recent edit
	result := make(map[int64]*generic.Update)

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
			editInfo := generic.TextMessageEditedContent{
				MessageID: messageID,
				NewText:   newText,
			}

			result[messageID] = &generic.Update{
				UpdateID:   updateID,
				ChatID:     uuid.UUID(chatID),
				SenderID:   senderID,
				UpdateType: domain.UpdateTypeTextMessageEdited,
				CreatedAt:  createdAt.Unix(),
				Content: generic.UpdateContent{
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
	updates []generic.Update,
) error {
	ids := updateTypesIDs(updates, domain.UpdateTypeTextMessageEdited)
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

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...)...)
	if err != nil {
		return err
	}
	defer rows.Close()

	edits := make(map[int64]generic.TextMessageEditedContent)

	for rows.Next() {
		var (
			updateID  int64
			newText   string
			messageID int64
		)

		if err := rows.Scan(&updateID, &newText, &messageID); err != nil {
			return err
		}

		edits[updateID] = generic.TextMessageEditedContent{
			NewText:   newText,
			MessageID: messageID,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the edit info
	for i, update := range updates {
		if update.UpdateType == domain.UpdateTypeTextMessageEdited {
			if info, ok := edits[update.UpdateID]; ok {
				updates[i].Content.TextMessageEdited = &info
			} else {
				panic("database is inconsistent!!!")
			}
		}
	}

	return nil
}

func (r *GenericUpdateRepository) fillFileMessages(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []generic.Update,
) error {
	ids := updateTypesIDs(updates, domain.UpdateTypeFileMessage)
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

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...)...)
	if err != nil {
		return err
	}
	defer rows.Close()

	fileMsgs := make(map[int64]generic.FileMessageContent)

	for rows.Next() {
		var (
			updateID      int64
			fileID        uuid.UUID
			fileName      string
			fileMimeType  string
			fileSize      int64
			fileURL       string
			fileCreatedAt int64
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

		fileMsgs[updateID] = generic.FileMessageContent{
			File: generic.FileMeta{
				FileId:    fileID,
				FileName:  fileName,
				MimeType:  fileMimeType,
				FileSize:  fileSize,
				FileURL:   fileURL,
				CreatedAt: fileCreatedAt,
			},
			ReplyTo: replyToID,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the file message info
	for i, update := range updates {
		if update.UpdateType == domain.UpdateTypeFileMessage {
			if info, ok := fileMsgs[update.UpdateID]; ok {
				reactions, err := r.getMessageReactions(ctx, db, chatID, domain.UpdateID(update.UpdateID))
				if err != nil {
					return err
				}
				if reactions != nil {
					info.Reactions = reactions
				}

				updates[i].Content.FileMessage = &info
			} else {
				panic("database is inconsistent!!!")
			}
		}
	}

	return nil
}

func (r *GenericUpdateRepository) getMessageReactions(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	messageID domain.UpdateID,
) ([]generic.Update, error) {
	q := `
	SELECT 
		u.update_id,
		u.created_at,
		u.sender_id,
		r.reaction
	FROM messaging.update u
		JOIN messaging.reaction_update r ON r.chat_id = u.chat_id AND r.update_id = u.update_id
	WHERE u.chat_id = $1 AND r.message_id = $2`

	rows, err := db.Query(ctx, q, chatID, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]generic.Update, 0)
	for rows.Next() {
		var (
			updateID     int64
			createdAt    time.Time
			senderID     uuid.UUID
			reactionType string
		)
		if err := rows.Scan(&updateID, &createdAt, &senderID, &reactionType); err != nil {
			return nil, err
		}

		res = append(res, generic.Update{
			UpdateID:   updateID,
			ChatID:     uuid.UUID(chatID),
			SenderID:   senderID,
			UpdateType: "reaction",
			CreatedAt:  createdAt.Unix(),
			Content: generic.UpdateContent{
				Reaction: &generic.ReactionContent{
					Reaction:  reactionType,
					MessageID: int64(messageID),
				},
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *GenericUpdateRepository) fillReactions(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []generic.Update,
) error {
	ids := updateTypesIDs(updates, domain.UpdateTypeReaction)
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

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...)...)
	if err != nil {
		return err
	}
	defer rows.Close()

	reactions := make(map[int64]generic.ReactionContent)

	for rows.Next() {
		var (
			updateID  int64
			reaction  string
			messageID int64
		)

		if err := rows.Scan(&updateID, &reaction, &messageID); err != nil {
			return err
		}

		reactions[updateID] = generic.ReactionContent{
			Reaction:  reaction,
			MessageID: messageID,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the reaction info
	for i, update := range updates {
		if update.UpdateType == domain.UpdateTypeReaction {
			if info, ok := reactions[update.UpdateID]; ok {
				updates[i].Content.Reaction = &info
			} else {
				panic("database is inconsistent!!!")
			}
		}
	}

	return nil
}

func (r *GenericUpdateRepository) fillUpdateDeleted(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []generic.Update,
) error {
	ids := updateTypesIDs(updates, domain.UpdateTypeDeleted)
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

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...)...)
	if err != nil {
		return err
	}
	defer rows.Close()

	deletedUpdates := make(map[int64]generic.DeletedContent)

	for rows.Next() {
		var (
			updateID        int64
			deletedUpdateID int64
			mode            string
		)

		if err := rows.Scan(&updateID, &deletedUpdateID, &mode); err != nil {
			return err
		}

		deletedUpdates[updateID] = generic.DeletedContent{
			DeletedID:   deletedUpdateID,
			DeletedMode: mode,
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the deleted update info
	for i, update := range updates {
		if update.UpdateType == domain.UpdateTypeDeleted {
			if info, ok := deletedUpdates[update.UpdateID]; ok {
				updates[i].Content.Deleted = &info
			} else {
				panic("database is inconsistent!!!")
			}
		}
	}

	return nil
}

func (r *GenericUpdateRepository) fillSecretUpdates(
	ctx context.Context,
	db storage.ExecQuerier,
	chatID domain.ChatID,
	updates []generic.Update,
) error {
	ids := updateTypesIDs(updates, domain.UpdateTypeSecret)
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

	rows, err := db.Query(ctx, q, append([]any{chatID}, idsToAny(ids)...)...)
	if err != nil {
		return err
	}
	defer rows.Close()

	secretUpdates := make(map[int64]generic.SecretUpdateContent)

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

		secretUpdates[updateID] = generic.SecretUpdateContent{
			PayloadBase64:              base64.StdEncoding.EncodeToString(payload),
			KeyHashBase64:              base64.StdEncoding.EncodeToString(keyHash),
			InitializationVectorBase64: base64.StdEncoding.EncodeToString(initializationVector),
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Update the updates with the secret update info
	for i, update := range updates {
		if update.UpdateType == domain.UpdateTypeSecret {
			if info, ok := secretUpdates[update.UpdateID]; ok {
				updates[i].Content.Secret = &info
			} else {
				panic("database is inconsistent!!!")
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

func updateTypesIDs(updates []generic.Update, typ string) []int64 {
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
