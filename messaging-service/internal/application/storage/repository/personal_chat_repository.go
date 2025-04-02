package repository

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/personal"
)

var (
	ErrNotFound = errors.New("not found")
)

//go:generate mockery
type PersonalChatRepository interface {
	// Should return ErrNotFound if entity is not found
	FindById(context.Context, storage.ExecQuerier, domain.ChatID) (*personal.PersonalChat, error)
	// Should return ErrNotFound if entity is not found.
	// Members order should NOT affect the result
	FindByMembers(context.Context, storage.ExecQuerier, [2]domain.UserID) (*personal.PersonalChat, error)
	Update(context.Context, storage.ExecQuerier, *personal.PersonalChat) (*personal.PersonalChat, error)
	Create(context.Context, storage.ExecQuerier, *personal.PersonalChat) (*personal.PersonalChat, error)
	Delete(context.Context, storage.ExecQuerier, domain.ChatID) error
}
