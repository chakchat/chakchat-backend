package services

import (
	"time"

	"github.com/google/uuid"
)

const (
	ChatTypePersonal       = "personal"
	ChatTypeGroup          = "group"
	ChatTypeSecretPersonal = "secret_personal"
	ChatTypeSecretGroup    = "secret_group"
)

type GenericChat struct {
	ChatID    uuid.UUID
	CreatedAt int64
	ChatType  string
	Members   []uuid.UUID
	Info      GenericChatInfo
}

// You should call one of methods depending on GenericChat.ChatType
type GenericChatInfo interface {
	// You should call it if chat type is `personal`. It will panic otherwise
	Personal() PersonalInfo
	// You should call it if chat type is `group`. It will panic otherwise
	Group() GroupInfo
	// You should call it if chat type is `secret_personal`. It will panic otherwise
	SecretPersonal() SecretPersonalInfo
	// You should call it if chat type is `secret_group`. It will panic otherwise
	SecretGroup() SecretGroupInfo
}

type PersonalInfo struct {
	BlockedBy []uuid.UUID
}

type GroupInfo struct {
	Admin            uuid.UUID
	GroupName        string
	GroupDescription string
	GroupPhoto       string
}

type SecretPersonalInfo struct {
	Expiration *time.Duration
}

type SecretGroupInfo struct {
	Admin            uuid.UUID
	GroupName        string
	GroupDescription string
	GroupPhoto       string
	Expiration       *time.Duration
}

func NewPersonalGenericChat(
	id uuid.UUID,
	createdAt int64,
	members []uuid.UUID,
	info PersonalInfo,
) GenericChat {
	return GenericChat{
		ChatID:    id,
		CreatedAt: createdAt,
		ChatType:  ChatTypePersonal,
		Members:   members,
		Info:      &genericChatInfo{personal: &info},
	}
}

func NewGroupGenericChat(
	id uuid.UUID,
	createdAt int64,
	members []uuid.UUID,
	info GroupInfo,
) GenericChat {
	return GenericChat{
		ChatID:    id,
		CreatedAt: createdAt,
		ChatType:  ChatTypeGroup,
		Members:   members,
		Info:      &genericChatInfo{group: &info},
	}
}

func NewSecretPersonalGenericChat(
	id uuid.UUID,
	createdAt int64,
	members []uuid.UUID,
	info SecretPersonalInfo,
) GenericChat {
	return GenericChat{
		ChatID:    id,
		CreatedAt: createdAt,
		ChatType:  ChatTypeSecretPersonal,
		Members:   members,
		Info:      &genericChatInfo{secretPersonal: &info},
	}
}

func NewSecretGroupGenericChat(
	id uuid.UUID,
	createdAt int64,
	members []uuid.UUID,
	info SecretGroupInfo,
) GenericChat {
	return GenericChat{
		ChatID:    id,
		CreatedAt: createdAt,
		ChatType:  ChatTypeSecretGroup,
		Members:   members,
		Info:      &genericChatInfo{secretGroup: &info},
	}
}

type genericChatInfo struct {
	personal       *PersonalInfo
	group          *GroupInfo
	secretPersonal *SecretPersonalInfo
	secretGroup    *SecretGroupInfo
}

func (i *genericChatInfo) Personal() PersonalInfo             { return *i.personal }
func (i *genericChatInfo) Group() GroupInfo                   { return *i.group }
func (i *genericChatInfo) SecretPersonal() SecretPersonalInfo { return *i.secretPersonal }
func (i *genericChatInfo) SecretGroup() SecretGroupInfo       { return *i.secretGroup }
