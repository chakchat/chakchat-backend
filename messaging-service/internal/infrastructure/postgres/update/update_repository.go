package update

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type UpdateRepository struct{}

func NewUpdateRepository() *UpdateRepository {
	return &UpdateRepository{}
}

func (r *UpdateRepository) FindGenericMessage(
	context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID,
) (*domain.Message, error)

func (r *UpdateRepository) DeleteMessage(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) error
func (r *UpdateRepository) CreateUpdateDeleted(context.Context, storage.ExecQuerier, *domain.UpdateDeleted) (*domain.UpdateDeleted, error)

func (r *UpdateRepository) CreateTextMessage(context.Context, storage.ExecQuerier, *domain.TextMessage) (*domain.TextMessage, error)
func (r *UpdateRepository) CreateTextMessageEdited(context.Context, storage.ExecQuerier, *domain.TextMessageEdited) (*domain.TextMessageEdited, error)
func (r *UpdateRepository) FindTextMessage(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) (*domain.TextMessage, error)
func (r *UpdateRepository) UpdateTextMessage(context.Context, storage.ExecQuerier, *domain.TextMessage) (*domain.TextMessage, error)

func (r *UpdateRepository) CreateReaction(context.Context, storage.ExecQuerier, *domain.Reaction) (*domain.Reaction, error)
func (r *UpdateRepository) FindReaction(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) (*domain.Reaction, error)
func (r *UpdateRepository) DeleteReaction(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) error

func (r *UpdateRepository) FindFileMessage(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) (*domain.FileMessage, error)
func (r *UpdateRepository) CreateFileMessage(context.Context, storage.ExecQuerier, *domain.FileMessage) (*domain.FileMessage, error)
