package services

import "errors"

var (
	ErrSignInKeyNotFound = errors.New("sign in key not found")
	ErrWrongCode         = errors.New("wrong phone verification code")
)

type JWT string

type TokenPair struct {
	Access  JWT
	Refresh JWT
}
