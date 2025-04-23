package domain

import (
	"time"
)

const (
	UpdateTypeTextMessage       = "text_message"
	UpdateTypeTextMessageEdited = "text_message_edited"
	UpdateTypeFileMessage       = "file_message"
	UpdateTypeReaction          = "reaction"
	UpdateTypeDeleted           = "update_deleted"
	UpdateTypeSecret            = "secret_update"
)

type (
	UpdateID uint64
)

type Timestamp int64

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

type Update struct {
	// It will be assigned automatically when it is stored in DB
	UpdateID UpdateID
	ChatID   ChatID
	SenderID UserID

	CreatedAt Timestamp
	Deleted   []*UpdateDeleted
}

type DeleteMode string

func NewDeleteMode(mode string) (DeleteMode, error) {
	switch mode {
	case DeleteModeForSender:
		return DeleteModeForSender, nil
	case DeleteModeForAll:
		return DeleteModeForAll, nil
	default:
		return "", ErrInvalidDeleteMode
	}
}

const (
	DeleteModeForSender = "for_deletion_sender"
	DeleteModeForAll    = "for_all"
)

// Caution: SenderID here is deletion update's sender but not a original update's sender
type UpdateDeleted struct {
	Update
	DeletedID UpdateID
	Mode      DeleteMode
}

func (u *Update) DeletedForAll() bool {
	for _, d := range u.Deleted {
		if d.Mode == DeleteModeForAll {
			return true
		}
	}

	return false
}

func (u *Update) DeletedFor(user UserID) bool {
	for _, d := range u.Deleted {
		if d.Mode == DeleteModeForSender && d.SenderID == user {
			return true
		}
		if d.Mode == DeleteModeForAll {
			return true
		}
	}

	return false
}

func (u *Update) AddDeletion(sender UserID, mode DeleteMode) {
	d := &UpdateDeleted{
		Update: Update{
			ChatID:   u.ChatID,
			SenderID: sender,
		},
		DeletedID: u.UpdateID,
		Mode:      mode,
	}

	if mode == DeleteModeForAll {
		u.Deleted = []*UpdateDeleted{d}
	} else {
		u.Deleted = append(u.Deleted, d)
	}
}
