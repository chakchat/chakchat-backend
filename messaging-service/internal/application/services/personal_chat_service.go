package services

import (
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrChatNotFound         = errors.New("service: chat not found")
	ErrChatAlreadyBlocked   = errors.New("service: chat already blocked")
	ErrChatAlreadyUnblocked = errors.New("service: chat already unblocked")
	ErrUserNotMember        = errors.New("service: user not chat member")

	ErrChatAlreadyExists = errors.New("service: chat already exists")
	ErrChatWithMyself    = errors.New("service: chat with myself")

	ErrUnknown = errors.New("service: unknown error")
)

type PersonalChatService struct {
	repo repository.PersonalChatRepository
}

func NewPersonalChatService(repo repository.PersonalChatRepository) *PersonalChatService {
	return &PersonalChatService{
		repo: repo,
	}
}

func (s *PersonalChatService) BlockChat(userId, chatId uuid.UUID) error {
	chat, err := s.repo.FindById(domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChatNotFound
		}
		return errors.Join(ErrUnknown, err)
	}

	err = chat.BlockBy(domain.UserID(userId))

	if err != nil {
		if errors.Is(err, domain.ErrAlreadyBlocked) {
			return ErrChatAlreadyBlocked
		}
		if errors.Is(err, domain.ErrUserNotMember) {
			return ErrUserNotMember
		}
		return errors.Join(ErrUnknown, err)
	}

	if _, err := s.repo.Update(chat); err != nil {
		return errors.Join(ErrUnknown, err)
	}

	return nil
}

func (s *PersonalChatService) UnblockChat(userId, chatId uuid.UUID) error {
	chat, err := s.repo.FindById(domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChatNotFound
		}
		return errors.Join(ErrUnknown, err)
	}

	err = chat.UnblockBy(domain.UserID(userId))

	if err != nil {
		if errors.Is(err, domain.ErrAlreadyUnblocked) {
			return ErrChatAlreadyUnblocked
		}
		if errors.Is(err, domain.ErrUserNotMember) {
			return ErrUserNotMember
		}
		return errors.Join(ErrUnknown, err)
	}

	if _, err := s.repo.Update(chat); err != nil {
		return errors.Join(ErrUnknown, err)
	}

	return nil
}

func (s *PersonalChatService) CreateChat(members [2]uuid.UUID) (*dto.PersonalChatDTO, error) {
	domainMembers := [2]domain.UserID{domain.UserID(members[0]), domain.UserID(members[1])}

	if err := s.validateChatNotExists(domainMembers); err != nil {
		return nil, err
	}

	chat, err := domain.NewPersonalChat(domainMembers)

	if err != nil {
		if errors.Is(err, domain.ErrChatWithMyself) {
			return nil, ErrChatWithMyself
		}
		return nil, errors.Join(ErrUnknown, err)
	}

	chat, err = s.repo.Create(chat)
	if err != nil {
		return nil, errors.Join(ErrUnknown, err)
	}

	chatDto := dto.NewPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *PersonalChatService) CreateSecretChat(members [2]uuid.UUID) (*dto.PersonalChatDTO, error) {
	domainMembers := [2]domain.UserID{domain.UserID(members[0]), domain.UserID(members[1])}

	if err := s.validateChatNotExists(domainMembers); err != nil {
		return nil, err
	}

	chat, err := domain.NewSecretPersonalChat(domainMembers)

	if err != nil {
		if errors.Is(err, domain.ErrChatWithMyself) {
			return nil, ErrChatWithMyself
		}
		return nil, errors.Join(ErrUnknown, err)
	}

	chat, err = s.repo.Create(chat)
	if err != nil {
		return nil, errors.Join(ErrUnknown, err)
	}

	chatDto := dto.NewPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *PersonalChatService) GetChatById(chatId uuid.UUID) (*dto.PersonalChatDTO, error) {
	chat, err := s.repo.FindById(domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, errors.Join(ErrUnknown, err)
	}

	chatDto := dto.NewPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *PersonalChatService) validateChatNotExists(members [2]domain.UserID) error {
	_, err := s.repo.FindByMembers(members)

	if err != nil && err != repository.ErrNotFound {
		return errors.Join(ErrUnknown, err)
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return ErrChatAlreadyExists
	}

	return nil
}
