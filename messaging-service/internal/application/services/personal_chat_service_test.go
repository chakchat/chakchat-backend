package services

import (
	"testing"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository/mocks"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPersonalChat_CreateChat(t *testing.T) {
	user1 := uuid.MustParse("97abd308-8878-41ce-ba61-0a75c836a908")
	user2 := uuid.MustParse("0e0bdfc6-3a45-4a39-bdfe-6d6145842454")

	t.Run("Success", func(t *testing.T) {
		repo := mocks.NewMockPersonalChatRepository(t)
		repo.EXPECT().
			FindByMembers(mock.Anything).
			Return(nil, repository.ErrNotFound)

		repo.EXPECT().
			Create(mock.Anything).
			RunAndReturn(func(chat *domain.PersonalChat) (*domain.PersonalChat, error) {
				return chat, nil
			})

		service := NewPersonalChatService(repo)

		chat, err := service.CreateChat([2]uuid.UUID{user1, user2})

		require.NoError(t, err)
		require.Equal(t, [2]uuid.UUID{user1, user2}, chat.Members)
		require.False(t, chat.Blocked)
		require.False(t, chat.Secret)

		repo.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("SecretSuccess", func(t *testing.T) {
		repo := mocks.NewMockPersonalChatRepository(t)
		repo.EXPECT().
			FindByMembers(mock.Anything).
			Return(nil, repository.ErrNotFound)

		repo.EXPECT().
			Create(mock.Anything).
			RunAndReturn(func(chat *domain.PersonalChat) (*domain.PersonalChat, error) {
				return chat, nil
			})

		service := NewPersonalChatService(repo)

		chat, err := service.CreateSecretChat([2]uuid.UUID{user1, user2})

		require.NoError(t, err)
		require.Equal(t, [2]uuid.UUID{user1, user2}, chat.Members)
		require.False(t, chat.Blocked)
		require.True(t, chat.Secret)

		repo.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		repo := mocks.NewMockPersonalChatRepository(t)
		repo.EXPECT().
			FindByMembers(mock.Anything).
			Return(&domain.PersonalChat{}, nil)

		service := NewPersonalChatService(repo)

		_, err := service.CreateChat([2]uuid.UUID{user1, user2})

		require.ErrorIs(t, err, ErrChatAlreadyExists)

		repo.AssertNotCalled(t, "Create")
	})
}
