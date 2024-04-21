package inmem_test

import (
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/storage/user"
	"github.com/dvid-messanger/internal/storage/user/inmem"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"strconv"
	"testing"
)

func TestUserStorage_Save(t *testing.T) {
	op := "TestUserStorage_Save"
	t.Parallel()

	for i := 1; i <= 10; i++ {
		t.Run(fmt.Sprintf("%s i = %d", op, i), func(t *testing.T) {
			t.Parallel()

			storage := inmem.New()

			expected := genUser("test-mail"+strconv.Itoa(i), "test-bio")

			actualR, err := storage.Save(context.Background(), expected.Email, expected.Bio)
			require.NoError(t, err, "save should not error")
			assert.NotEmpty(t, actualR.Id, "returned id not empty")

			expected.Id = actualR.Id
			assert.Equal(t, *expected, actualR, "returned user match expected")

			saved, err := storage.User(context.Background(), expected.Id)
			require.NoError(t, err, "user should not error")
			assert.Equal(t, *expected, saved)
		})
	}

	t.Run(fmt.Sprintf("%s_EmailExistsReturnError", op), func(t *testing.T) {
		t.Parallel()

		storage := inmem.New()

		user1 := genUser("test-mail1", "test-bio1")
		user2 := genUser("test-mail2", "test-bio2")

		_, err := storage.Save(context.Background(), user1.Email, user1.Bio)
		require.NoError(t, err, "save should not error")

		_, err = storage.Save(context.Background(), user1.Email, user2.Bio)
		require.Error(t, err, "save with existing id should error")
		assert.ErrorIs(t, err, user.ErrUserExists, "error type correct")
	})
}

func genUser(email string, bio string) *model.User {
	uid := [16]byte(uuid.New())
	return &model.User{
		Id:    uid[:],
		Email: email,
		Bio:   bio,
	}
}
