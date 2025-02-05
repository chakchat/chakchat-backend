package domain

import "time"

const NoExpiration time.Duration = 0

var TimeFunc = func() time.Time {
	return time.Now()
}

type SecretChat struct {
	Chat

	// If no expiration is set it is equal to NoExpiration
	Expiration time.Duration
}

type SecretChatter interface {
	Chatter
	// If no expiration is set it is equal to NoExpiration
	Expiration() time.Duration
}
