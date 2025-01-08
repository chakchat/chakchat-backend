package services

import "errors"

var (
	ErrUsernameAlreadyExists     = errors.New("username already exists")
	ErrCreateUserValidationError = errors.New("create user validation error")
)

type CreateUserData struct {
	Username string
	Name     string
}
