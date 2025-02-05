package query

import "github.com/google/uuid"

type CreateGroupRequest struct {
	Admin   uuid.UUID
	Members []uuid.UUID
	Name    string
}

type UpdateGroupInfoRequest struct {
	ChatID      uuid.UUID
	Name        string
	Description string
}

type CreateSecretGroupRequest struct {
	Admin   uuid.UUID
	Members []uuid.UUID
	Name    string
}

type UpdateSecretGroupInfoRequest struct {
	ChatID      uuid.UUID
	Name        string
	Description string
}
