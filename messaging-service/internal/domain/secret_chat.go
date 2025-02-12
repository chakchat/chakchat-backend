package domain

import "time"

type SecretChat struct {
	Chat

	Exp *time.Duration
}

func (c *SecretChat) Expiration() *time.Duration {
	return c.Exp
}

type SecretChatter interface {
	Chatter
	Expiration() *time.Duration
}
