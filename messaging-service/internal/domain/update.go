package domain

import "errors"

type (
	UpdateID  uint64
	Timestamp int64
)

var (
	ErrUserNotSender     = errors.New("user is not update's sender")
	ErrUpdateNotFromChat = errors.New("update is not from this chat")
)

type Update struct {
	// It will be assigned automatically when it is stored in DB
	UpdateID UpdateID
	ChatID   ChatID
	SenderID UserID

	CreatedAt Timestamp
	Deleted   *UpdateDeleted
}

type DeleteMode int

const (
	DeleteModeForSender = iota
	DeleteModeForAll
)

type UpdateDeleted struct {
	Update
	DeletedID UpdateID
	Mode      DeleteMode
}
