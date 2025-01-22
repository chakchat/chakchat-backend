package entities

type GroupChat struct {
	ID      ChatID
	Admin   UserID
	Members []UserID

	Secret      bool
	Name        string
	Description string
	GroupPhoto  URL
	CreatedAt   Timestamp
}
