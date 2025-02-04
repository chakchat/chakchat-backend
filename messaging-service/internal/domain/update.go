package domain

import "errors"

type (
	UpdateID  uint64
	Timestamp int64
)

var (
	ErrUserNotSender     = errors.New("user is not update's sender")
	ErrUpdateNotFromChat = errors.New("update is not from this chat")
	ErrUpdateDeleted     = errors.New("update is deleted")
)

type Update struct {
	// It will be assigned automatically when it is stored in DB
	UpdateID UpdateID
	ChatID   ChatID
	SenderID UserID

	CreatedAt Timestamp
	Deleted   []UpdateDeleted
}

type DeleteMode int

const (
	DeleteModeForSender = iota
	DeleteModeForAll
)

// Caution: SenderID here is deletion update's sender but not a original update's sender
type UpdateDeleted struct {
	Update
	DeletedID UpdateID
	Mode      DeleteMode
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
	d := UpdateDeleted{
		Update: Update{
			ChatID:   u.ChatID,
			SenderID: sender,
		},
		DeletedID: u.UpdateID,
		Mode:      mode,
	}

	if mode == DeleteModeForAll {
		u.Deleted = []UpdateDeleted{d}
	} else {
		u.Deleted = append(u.Deleted, d)
	}
}
