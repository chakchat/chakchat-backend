package request

import "github.com/google/uuid"

type CreateGroup struct {
	SenderID uuid.UUID
	Members  []uuid.UUID
	Name     string
}

type DeleteGroup struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
}

type AddMember struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	MemberID uuid.UUID
}

type DeleteMember struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	MemberID uuid.UUID
}

type UpdateGroupPhoto struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	FileID   uuid.UUID
}

type DeleteGroupPhoto struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
}

type UpdateGroupInfo struct {
	ChatID      uuid.UUID
	SenderID    uuid.UUID
	Name        string
	Description string
}

type CreateSecretGroup struct {
	SenderID uuid.UUID
	Members  []uuid.UUID
	Name     string
}

type UpdateSecretGroupInfo struct {
	ChatID      uuid.UUID
	SenderID    uuid.UUID
	Name        string
	Description string
}
