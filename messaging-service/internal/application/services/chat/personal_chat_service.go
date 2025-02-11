package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/personal"
	"github.com/google/uuid"
)

var (
	ErrChatNotFound         = errors.New("service: chat not found")
	ErrChatAlreadyBlocked   = errors.New("service: chat already blocked")
	ErrChatAlreadyUnblocked = errors.New("service: chat already unblocked")
	ErrUserNotMember        = errors.New("service: user not chat member")

	ErrChatAlreadyExists = errors.New("service: chat already exists")
	ErrChatWithMyself    = errors.New("service: chat with myself")

	ErrInternal = errors.New("service: unknown error")
)

type PersonalChatService struct {
	repo repository.PersonalChatRepository
	pub  publish.Publisher
}

func NewPersonalChatService(repo repository.PersonalChatRepository, pub publish.Publisher) *PersonalChatService {
	return &PersonalChatService{
		repo: repo,
		pub:  pub,
	}
}

func (s *PersonalChatService) BlockChat(ctx context.Context, userId, chatId uuid.UUID) error {
	chat, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChatNotFound
		}
		return errors.Join(ErrInternal, err)
	}

	err = chat.BlockBy(domain.UserID(userId))

	if err != nil {
		if errors.Is(err, personal.ErrAlreadyBlocked) {
			return ErrChatAlreadyBlocked
		}
		if errors.Is(err, domain.ErrUserNotMember) {
			return ErrUserNotMember
		}
		return errors.Join(ErrInternal, err)
	}

	if _, err := s.repo.Update(ctx, chat); err != nil {
		return errors.Join(ErrInternal, err)
	}

	s.pub.PublishForUsers(getSecondUserSlice(chat.Members, domain.UserID(userId)), events.ChatBlocked{
		ChatID: chatId,
	})

	return nil
}

func (s *PersonalChatService) UnblockChat(ctx context.Context, userId, chatId uuid.UUID) error {
	chat, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChatNotFound
		}
		return errors.Join(ErrInternal, err)
	}

	err = chat.UnblockBy(domain.UserID(userId))

	if err != nil {
		if errors.Is(err, personal.ErrAlreadyUnblocked) {
			return ErrChatAlreadyUnblocked
		}
		if errors.Is(err, domain.ErrUserNotMember) {
			return ErrUserNotMember
		}
		return errors.Join(ErrInternal, err)
	}

	if _, err := s.repo.Update(ctx, chat); err != nil {
		return errors.Join(ErrInternal, err)
	}

	s.pub.PublishForUsers(getSecondUserSlice(chat.Members, domain.UserID(userId)), events.ChatUnblocked{
		ChatID: chatId,
	})

	return nil
}

func (s *PersonalChatService) CreateChat(ctx context.Context, userId, withUserId uuid.UUID) (*dto.PersonalChatDTO, error) {
	domainMembers := [2]domain.UserID{domain.UserID(userId), domain.UserID(withUserId)}

	if err := s.validateChatNotExists(ctx, domainMembers); err != nil {
		return nil, err
	}

	chat, err := personal.NewPersonalChat(domainMembers)

	if err != nil {
		if errors.Is(err, domain.ErrChatWithMyself) {
			return nil, ErrChatWithMyself
		}
		return nil, errors.Join(ErrInternal, err)
	}

	chat, err = s.repo.Create(ctx, chat)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	chatDto := dto.NewPersonalChatDTO(chat)

	s.pub.PublishForUsers([]uuid.UUID{withUserId}, events.ChatCreated{
		ChatID:   chatDto.ID,
		ChatType: events.ChatTypePersonal,
	})

	return &chatDto, nil
}

func (s *PersonalChatService) GetChatById(ctx context.Context,
	chatId uuid.UUID) (*dto.PersonalChatDTO, error) {
	chat, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, errors.Join(ErrInternal, err)
	}

	chatDto := dto.NewPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *PersonalChatService) DeleteChat(ctx context.Context, chatId uuid.UUID, deleteForAll bool) error {
	chat, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChatNotFound
		}
		return errors.Join(ErrInternal, err)
	}

	// TODO: put other logic here after you decide what to do with messages
	// For now I think that messages with in deleted chat should be deleted by background task

	if err := s.repo.Delete(ctx, chat.ID); err != nil {
		return errors.Join(ErrInternal, err)
	}

	s.pub.PublishForUsers(dto.UUIDs(chat.Members[:]), events.ChatDeleted{
		ChatID: chatId,
	})

	return nil
}

func (s *PersonalChatService) validateChatNotExists(ctx context.Context, members [2]domain.UserID) error {
	_, err := s.repo.FindByMembers(ctx, members)

	if err != nil && err != repository.ErrNotFound {
		return errors.Join(ErrInternal, err)
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return ErrChatAlreadyExists
	}

	return nil
}

func getSecondUserSlice(users [2]domain.UserID, first domain.UserID) []uuid.UUID {
	var second domain.UserID
	if users[0] == first {
		second = users[1]
	} else {
		second = users[0]
	}
	return []uuid.UUID{uuid.UUID(second)}
}
