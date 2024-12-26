package services

import "errors"

var (
	ErrRefreshTokenExpired     = errors.New("refresh token expired")
	ErrRefreshTokenInvalidated = errors.New("refresh token invalidated")
	ErrInvalidJWT              = errors.New("jwt token is invalid")
)
