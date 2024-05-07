package user

import (
	"bytes"
	"context"
	"errors"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/driver/secondary/storage/user"
	"github.com/dvid-messanger/test/mocks/mock_user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

var log = slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

func TestCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockUserSaver := &mock_user.MockUserSaver{}
		mockUserProvider := &mock_user.MockUserProvider{}

		expectedUser := model.User{
			Id:    []byte("userId"),
			Email: "test@example.com",
			Bio:   "This is a test",
		}

		mockUserSaver.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(expectedUser, nil)

		service := NewUser(log, mockUserSaver, mockUserProvider)

		usr, err := service.Create(context.Background(), "test@example.com", "This is a test")
		require.NoError(t, err)
		assert.Equal(t, expectedUser, *usr)
	})
	t.Run("UserExistsError", func(t *testing.T) {
		t.Parallel()

		mockUserSaver := &mock_user.MockUserSaver{}
		mockUserProvider := &mock_user.MockUserProvider{}

		mockUserSaver.On("Save", mock.Anything, mock.Anything, mock.Anything).
			Return(model.User{}, user.ErrUserExists)

		service := NewUser(log, mockUserSaver, mockUserProvider)

		_, err := service.Create(context.Background(), "test@example.com", "This is a test")
		assert.ErrorIs(t, err, ErrUserExists)
	})
	t.Run("UserSaveFailedError", func(t *testing.T) {
		t.Parallel()

		mockUserSaver := &mock_user.MockUserSaver{}
		mockUserProvider := &mock_user.MockUserProvider{}

		mockUserSaver.On("Save", mock.Anything, mock.Anything, mock.Anything).
			Return(model.User{}, errors.New("failed to save user"))

		service := NewUser(log, mockUserSaver, mockUserProvider)

		_, err := service.Create(context.Background(), "test@example.com", "This is a test")
		assert.Error(t, err)
	})
}

func TestUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockUserSaver := &mock_user.MockUserSaver{}
		mockUserProvider := &mock_user.MockUserProvider{}

		expectedUser := model.User{
			Id:    []byte("userId"),
			Email: "test@example.com",
			Bio:   "This is a test",
		}

		mockUserProvider.On("User", mock.Anything, mock.Anything).Return(expectedUser, nil)

		service := NewUser(log, mockUserSaver, mockUserProvider)

		usr, err := service.User(context.Background(), []byte("userID"))
		require.NoError(t, err)
		assert.Equal(t, expectedUser, *usr)
	})
	t.Run("UserNotFoundError", func(t *testing.T) {
		t.Parallel()

		mockUserSaver := &mock_user.MockUserSaver{}
		mockUserProvider := &mock_user.MockUserProvider{}

		mockUserProvider.On("User", mock.Anything, mock.Anything).
			Return(model.User{}, user.ErrUserNotFound)

		service := NewUser(log, mockUserSaver, mockUserProvider)

		_, err := service.User(context.Background(), []byte("nonexistentUserID"))
		assert.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestUsers(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockUserSaver := &mock_user.MockUserSaver{}
		mockUserProvider := &mock_user.MockUserProvider{}

		expectedUsers := []model.User{
			{Id: []byte("user1"), Email: "user1@example.com", Bio: "User 1"},
			{Id: []byte("user2"), Email: "user2@example.com", Bio: "User 2"},
		}

		mockUserProvider.On("Users", mock.Anything).Return(expectedUsers, nil)

		service := NewUser(log, mockUserSaver, mockUserProvider)

		users, err := service.Users(context.Background())
		require.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
	})
	t.Run("NoUsersSuccess", func(t *testing.T) {
		t.Parallel()

		mockUserSaver := &mock_user.MockUserSaver{}
		mockUserProvider := &mock_user.MockUserProvider{}

		var expectedUsers []model.User

		mockUserProvider.On("Users", mock.Anything).Return(expectedUsers, nil)

		service := NewUser(log, mockUserSaver, mockUserProvider)

		users, err := service.Users(context.Background())
		require.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
	})
}
