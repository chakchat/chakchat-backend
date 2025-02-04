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

type UpdateDeletion struct {
	Update
	DeletedID UpdateID
	Mode      DeleteMode
}

func NewUpdateDeletion(u *Update, mode DeleteMode) UpdateDeletion {
	return UpdateDeletion{
		DeletedID: u.UpdateID,
		Mode:      mode,
	}
}
