package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/jackc/pgx/v5"
)

type ChatterRepository struct {
	// I am sorry for this cringe, just don't wanna copy-paste code
	personalRepo *PersonalChatRepository
	groupRepo    *GroupChatRepository
	secpPersonalRepo *SecretPersonalChatRepository
	secGroupRepo *SecretGroupChatRepository
}

func NewChatterRepository() *ChatterRepository {
	return &ChatterRepository{
		personalRepo:     NewPersonalChatRepository(),
		groupRepo:        NewGroupChatRepository(),
		secpPersonalRepo: NewSecretPersonalChatRepository(),
		secGroupRepo:     NewSecretGroupChatRepository(),
	}
}

func (r *ChatterRepository) FindChatter(
	ctx context.Context, db storage.ExecQuerier, id domain.ChatID,
) (domain.Chatter, error) {
	q := `SELECT chat_type FROM messaging.chat WHERE chat_id = $1`

	var (
		chatType string
	)
	row := db.QueryRow(ctx, q, id)

	err := row.Scan(&chatType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	if chatType == services.ChatTypePersonal {
		return r.personalRepo.FindById(ctx, db, id)
	}

	if chatType == services.ChatTypeGroup {
		return r.groupRepo.FindById(ctx, db, id)
	}

	if chatType == services.ChatTypeSecretPersonal {
		return r.secpPersonalRepo.FindById(ctx, db, id)
	}

	if chatType == services.ChatTypeSecretGroup {
		return r.secGroupRepo.FindById(ctx, db, id)
	}

	return nil, errors.Join(repository.ErrNotFound, fmt.Errorf("unknown Chatter type: %s", chatType))
}
