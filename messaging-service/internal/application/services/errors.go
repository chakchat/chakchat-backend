package services

import "errors"

var (
	ErrInternal          = errors.New("service: unknown error")
	ErrFileNotFound      = errors.New("service: file not found")
	ErrChatNotFound      = errors.New("service: chat not found")
	ErrInvalidPhoto      = errors.New("service: invalid photo")
	ErrChatAlreadyExists = errors.New("service: chat already exists")
	ErrMessageNotFound   = errors.New("service: message not found")
)
