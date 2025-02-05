package domain

import "time"

type SecretChat struct {
	Chat

	// If no expiration is set it is equal to NoExpiration
	Expiration *time.Duration
}

type SecretChatter interface {
	Chatter
	// If no expiration is set it is equal to NoExpiration
	Expiration() *time.Duration
}
