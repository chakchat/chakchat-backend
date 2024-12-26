package services

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrSendCodeFreqExceeded = errors.New("send code operation frequency exceeded")
)
