package auth

import (
	"bytes"
	"context"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/storage/auth"
	"github.com/dvid-messanger/mocks/mock_auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"testing"
	"time"
)

var log = slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

func TestLogin(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		passHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		require.NoError(t, err)

		mockUserProvider := &mock_auth.MockUserProvider{}
		mockUserProvider.On("User", mock.Anything, mock.Anything).Return(
			model.UserCredentials{
				Id:       []byte("mockUserID"),
				Email:    "test@example.com",
				PassHash: passHash,
			}, nil)

		mockTokenMaker := &mock_auth.MockTokenMaker{}
		mockTokenMaker.On("MakeToken", mock.Anything, mock.Anything).Return("mockToken", nil)

		service := NewService(log, nil, mockUserProvider, mockTokenMaker, time.Hour)

		token, err := service.Login(context.Background(), "test@example.com", "password")
		require.NoError(t, err)
		assert.Equal(t, "mockToken", token)
	})
	t.Run("InvalidCredentialsError", func(t *testing.T) {
		t.Parallel()

		mockUserProvider := &mock_auth.MockUserProvider{}
		mockUserProvider.On("User", mock.Anything, "test@example.com").Return(
			model.UserCredentials{}, auth.ErrUserNotFound)

		service := NewService(log, nil, mockUserProvider, nil, time.Hour)

		_, err := service.Login(context.Background(), "test@example.com", "wrongpassword")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		passHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		require.NoError(t, err)

		mockUserSaver := &mock_auth.MockUserSaver{}
		mockUserSaver.On("Save", mock.Anything, mock.Anything, "newuser@example.com", mock.Anything).Return(
			model.UserCredentials{
				Id:       []byte("mockUserID"),
				Email:    "newuser@example.com",
				PassHash: passHash, // password: "password"
			}, nil)

		service := NewService(log, mockUserSaver, nil, nil, time.Hour)

		uid := []byte("mockUserID")
		userID, err := service.Create(context.Background(), uid, "newuser@example.com", "password")
		require.NoError(t, err)
		assert.Equal(t, uid, userID)
	})
	t.Run("UserExistsError", func(t *testing.T) {
		t.Parallel()

		mockUserSaver := &mock_auth.MockUserSaver{}
		mockUserSaver.On("Save", mock.Anything, mock.Anything, "test@example.com", mock.Anything).Return(
			model.UserCredentials{}, auth.ErrUserExists)

		service := NewService(log, mockUserSaver, nil, nil, time.Hour)

		uid := []byte("mockUserID")
		_, err := service.Create(context.Background(), uid, "test@example.com", "password")
		assert.ErrorIs(t, err, ErrUserExists)
	})
}
