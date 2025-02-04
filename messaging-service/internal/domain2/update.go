package domain

type Update struct {
	// It will be assigned automatically when it is stored in DB
	UpdateID UpdateID
	ChatID   ChatID
	SenderID UserID

	SentAt Timestamp
}

type UpdateDeletion struct {
	Update
	DeletedID UpdateID
}
