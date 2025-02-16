package services

type Error struct {
	text string
}

func (e Error) Error() string {
	return e.text
}

var (
	ErrFileNotFound         = Error{"service: file not found"}
	ErrChatNotFound         = Error{"service: chat not found"}
	ErrInvalidPhoto         = Error{"service: invalid photo"}
	ErrChatAlreadyExists    = Error{"service: chat already exists"}
	ErrMessageNotFound      = Error{"service: message not found"}
	ErrReactionNotFound     = Error{"service: reaction not found"}
	ErrSecretUpdateNotFound = Error{"service: secret update is not found"}
)
