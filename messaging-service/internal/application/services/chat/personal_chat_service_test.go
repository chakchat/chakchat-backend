package chat

import (
	"context"
	"testing"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository/mocks"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/personal"
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
			FindByMembers(mock.Anything, mock.Anything).
			Return(nil, repository.ErrNotFound)

		repo.EXPECT().
			Create(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, chat *personal.PersonalChat) (*personal.PersonalChat, error) {
				return chat, nil
			})

		service := NewPersonalChatService(repo, FakePublisher{})

		chat, err := service.CreateChat(context.Background(), user1, user2)

		require.NoError(t, err)
		require.Equal(t, [2]uuid.UUID{user1, user2}, chat.Members)
		require.False(t, chat.Blocked)

		repo.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		repo := mocks.NewMockPersonalChatRepository(t)
		repo.EXPECT().
			FindByMembers(mock.Anything, mock.Anything).
			Return(&personal.PersonalChat{}, nil)

		service := NewPersonalChatService(repo, FakePublisher{})

		_, err := service.CreateChat(context.Background(), user1, user2)

		require.ErrorIs(t, err, ErrChatAlreadyExists)

		repo.AssertNotCalled(t, "Create")
	})
}

type FakePublisher struct{}

func (FakePublisher) PublishForUsers([]uuid.UUID, publish.Event) {}
