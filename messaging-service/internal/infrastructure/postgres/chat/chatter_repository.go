package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	ChatTypePersonal       = "personal"
	ChatTypeGroup          = "group"
	ChatTypeSecretPersonal = "secret_personal"
	ChatTypeSecretGroup    = "secret_group"
)

type ChatterRepository struct {
	db *pgx.Conn
	// I am sorry for this cringe, just don't wanna copy-paste code
	personalRepo *PersonalChatRepository
	groupRepo    *GroupChatRepository
}

func NewChatterRepository(db *pgx.Conn) *ChatterRepository {
	return &ChatterRepository{
		db:           db,
		personalRepo: NewPersonalChatRepository(db),
		groupRepo:    NewGroupChatRepository(db),
	}
}

func (r *ChatterRepository) FindChatter(ctx context.Context, id domain.ChatID) (domain.Chatter, error) {
	q := `SELECT chat_type FROM messaging.chat WHERE chat_id = $1`

	var (
		chatType string
	)
	row := r.db.QueryRow(ctx, q, id)

	err := row.Scan(&chatType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	if chatType == ChatTypePersonal {
		return r.personalRepo.FindById(ctx, id)
	}

	if chatType == ChatTypeGroup {
		return r.groupRepo.FindById(ctx, id)
	}

	return nil, errors.Join(repository.ErrNotFound, fmt.Errorf("unknown Chatter type: %s", chatType))
}
