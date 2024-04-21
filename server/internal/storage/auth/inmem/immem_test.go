package inmem_test

import (
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/storage/auth"
	"github.com/dvid-messanger/internal/storage/auth/inmem"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"strconv"
	"testing"
)

func TestAuthStorage_Save(t *testing.T) {
	op := "TestAuthStorage_Save"
	t.Parallel()

	for i := 1; i <= 10; i++ {
		t.Run(fmt.Sprintf("%s i = %d", op, i), func(t *testing.T) {
			t.Parallel()

			storage := inmem.New()
			expected := genCreds("test-mail"+strconv.Itoa(i), "test-password")

			actualR, err := storage.Save(context.Background(), expected.Id, expected.Email, expected.PassHash)
			require.NoError(t, err, "save should not error")

			assert.Equal(t, *expected, actualR, "returned creds match expected")

			user, err := storage.User(context.Background(), expected.Email)
			require.NoError(t, err, "user should not error")
			assert.Equal(t, *expected, user, "saved creds match expected")
		})
	}

	t.Run(fmt.Sprintf("%s_EmailExistsReturnError", op), func(t *testing.T) {
		t.Parallel()

		storage := inmem.New()

		creds1 := genCreds("test-mail1", "test-password1")
		creds2 := genCreds("test-mail2", "test-password2")

		_, err := storage.Save(context.Background(), creds1.Id, creds1.Email, creds1.PassHash)
		require.NoError(t, err, "save should not error")

		_, err = storage.Save(context.Background(), creds2.Id, creds1.Email, creds2.PassHash)
		require.Error(t, err, "save existing same email should error")
		assert.ErrorIs(t, err, auth.ErrUserExists, "error type correct")
	})

	t.Run(fmt.Sprintf("%s_IdExistsReturnError", op), func(t *testing.T) {
		t.Parallel()

		storage := inmem.New()

		creds1 := genCreds("test-mail1", "test-password1")
		creds2 := genCreds("test-mail2", "test-password2")

		_, err := storage.Save(context.Background(), creds1.Id, creds1.Email, creds1.PassHash)
		require.NoError(t, err, "save should not error")

		_, err = storage.Save(context.Background(), creds1.Id, creds2.Email, creds2.PassHash)
		require.Error(t, err, "save with existing id should error")
		assert.ErrorIs(t, err, auth.ErrUserExists, "error type correct")
	})
}

func genCreds(email string, pass string) *model.UserCredentials {
	uid := [16]byte(uuid.New())
	return &model.UserCredentials{
		Id:       uid[:],
		Email:    email,
		PassHash: []byte(pass),
	}
}
