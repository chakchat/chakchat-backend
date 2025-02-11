package request

import "github.com/google/uuid"

type CreateGroup struct {
	Admin   uuid.UUID
	Members []uuid.UUID
	Name    string
}

type UpdateGroupInfo struct {
	ChatID      uuid.UUID
	Name        string
	Description string
}

type CreateSecretGroup struct {
	Admin   uuid.UUID
	Members []uuid.UUID
	Name    string
}

type UpdateSecretGroupInfo struct {
	ChatID      uuid.UUID
	Name        string
	Description string
}
