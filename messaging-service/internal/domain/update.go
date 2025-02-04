package domain

type (
	UpdateID  uint64
	Timestamp int64
)

type Update struct {
	// It will be assigned automatically when it is stored in DB
	UpdateID UpdateID
	ChatID   ChatID
	SenderID UserID

	SentAt Timestamp
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

func NewUpdateDeleted(u *Update, mode DeleteMode) UpdateDeleted {
	return UpdateDeleted{
		DeletedID: u.UpdateID,
		Mode:      mode,
	}
}
