package inmem

import (
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/id"
	"github.com/dvid-messanger/internal/storage/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"strconv"
	"testing"
)

func TestStorage_SaveUser(t *testing.T) {
	storage := New()

	for i := 1; i <= 10; i++ {
		expected := genUser("test-mail"+strconv.Itoa(i), "test-bio")

		actualR, err := storage.SaveUser(context.Background(), expected.Email, expected.Bio)
		require.NoError(t, err, "save should not error")

		assert.NotEmpty(t, actualR.Id, "returned id not empty")
		assert.Equal(t, expected.Email, actualR.Email, "returned email match expected")
		assert.Equal(t, expected.Bio, actualR.Bio, "returned bio match expected")

		assert.Len(t, storage.hash, i, "hash len correct")

		actualS, ok := storage.hash[id.Id(actualR.Id)]
		assert.True(t, ok, "hash contains user")
		assert.Equal(t, actualR.Id, actualS.Id, "hash id match returned")
		assert.Equal(t, expected.Email, actualR.Email, "hash email match expected")
		assert.Equal(t, expected.Bio, actualS.Bio, "hash bio match expected")
	}
}

func TestStorage_SaveUser_SameEmailError(t *testing.T) {
	storage := New()

	user1 := genUser("test-mail", "test-bio")
	_, err := storage.SaveUser(context.Background(), user1.Email, user1.Bio)
	require.NoError(t, err, "save should not error")

	user2 := genUser("test-mail", "test-bio2")
	_, err = storage.SaveUser(context.Background(), user2.Email, user2.Bio)
	if assert.Error(t, err, "save with same email should error") {
		assert.ErrorIs(t, err, user.ErrUserExists, "error type correct")
	}
}

func TestStorage_User(t *testing.T) {
	storage := New()

	user1 := genUser("test-mail", "test-bio")
	user2 := genUser("test-mail2", "test-bio2")
	saved1, err := storage.SaveUser(context.Background(), user1.Email, user1.Bio)
	require.NoError(t, err, "save should not error")
	saved2, err := storage.SaveUser(context.Background(), user2.Email, user2.Bio)
	require.NoError(t, err, "save should not error")

	actual1, err := storage.User(context.Background(), saved1.Id)
	require.NoError(t, err, "user should not error")
	actual2, err := storage.User(context.Background(), saved2.Id)
	require.NoError(t, err, "user should not error")

	assert.Equal(t, saved1.Id, actual1.Id, "id match saved")
	assert.Equal(t, user1.Email, actual1.Email, "email match expected")
	assert.Equal(t, user1.Bio, actual1.Bio, "bio match expected")

	assert.Equal(t, saved2.Id, actual2.Id, "id match saved")
	assert.Equal(t, user2.Email, actual2.Email, "email match expected")
	assert.Equal(t, user2.Bio, actual2.Bio, "bio match expected")
}

func genUser(email string, bio string) *model.User {
	uid := [16]byte(uuid.New())
	return &model.User{
		Id:    uid[:],
		Email: email,
		Bio:   bio,
	}
}
