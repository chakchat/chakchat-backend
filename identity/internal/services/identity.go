package services

import "errors"

var (
	ErrAccessTokenExpired = errors.New("access token expired")
)
