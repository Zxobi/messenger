package inmem

import (
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/storage/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"strconv"
	"testing"
)

func TestAuthStorage_Save(t *testing.T) {
	storage := New()

	for i := 1; i <= 10; i++ {
		expected := genCreds("test-mail"+strconv.Itoa(i), "test-password")

		actualR, err := storage.Save(context.Background(), expected.Id, expected.Email, expected.PassHash)
		require.NoError(t, err, "save should not error")

		assert.ElementsMatch(t, expected.Id, actualR.Id, "returned id match expected")
		assert.Equal(t, expected.Email, actualR.Email, "returned email match expected")
		assert.Equal(t, expected.PassHash, actualR.PassHash, "returned passHash match expected")

		assert.Len(t, storage.hash, i, "hash len correct")

		actualS, ok := storage.hash[expected.Email]
		assert.True(t, ok, "hash contains creds")
		assert.ElementsMatch(t, expected.Id, actualS.Id, "hash id match expected")
		assert.Equal(t, expected.Email, actualS.Email, "hash email match expected")
		assert.Equal(t, expected.PassHash, actualS.PassHash, "hash passHash match expected")
	}
}

func TestAuthStorage_Save_SameEmailError(t *testing.T) {
	storage := New()

	creds1 := genCreds("test-mail", "test-password")
	_, err := storage.Save(context.Background(), creds1.Id, creds1.Email, creds1.PassHash)
	require.NoError(t, err, "save should not error")

	creds2 := genCreds("test-mail", "test-password2")
	_, err = storage.Save(context.Background(), creds2.Id, creds1.Email, creds2.PassHash)
	if assert.Error(t, err, "save with same email should error") {
		assert.ErrorIs(t, err, auth.ErrUserExists, "error type correct")
	}
}

func TestAuthStorage_Save_SameIdError(t *testing.T) {
	storage := New()

	creds1 := genCreds("test-mail", "test-password")
	_, err := storage.Save(context.Background(), creds1.Id, creds1.Email, creds1.PassHash)
	require.NoError(t, err, "save should not error")

	creds2 := genCreds("test-mail2", "test-password2")
	_, err = storage.Save(context.Background(), creds1.Id, creds2.Email, creds2.PassHash)
	if assert.Error(t, err, "save with same id should error") {
		assert.ErrorIs(t, err, auth.ErrUserExists, "error type correct")
	}
}

func TestAuthStorage_User(t *testing.T) {
	storage := New()

	creds1 := genCreds("test-mail", "test-password")
	creds2 := genCreds("test-mail2", "test-password2")
	_, err := storage.Save(context.Background(), creds1.Id, creds1.Email, creds1.PassHash)
	require.NoError(t, err, "save should not error")
	_, err = storage.Save(context.Background(), creds2.Id, creds2.Email, creds2.PassHash)
	require.NoError(t, err, "save should not error")

	actual1, err := storage.User(context.Background(), creds1.Email)
	require.NoError(t, err, "user should not error")
	actual2, err := storage.User(context.Background(), creds2.Email)
	require.NoError(t, err, "user should not error")

	assert.ElementsMatch(t, creds1.Id, actual1.Id, "id match expected")
	assert.Equal(t, creds1.Email, actual1.Email, "email match expected")
	assert.Equal(t, creds1.PassHash, actual1.PassHash, "passHash match expected")

	assert.ElementsMatch(t, creds2.Id, actual2.Id, "id match expected")
	assert.Equal(t, creds2.Email, actual2.Email, "email match expected")
	assert.Equal(t, creds2.PassHash, actual2.PassHash, "passHash match expected")
}

func genCreds(email string, pass string) *model.UserCredentials {
	uid := [16]byte(uuid.New())
	return &model.UserCredentials{
		Id:       uid[:],
		Email:    email,
		PassHash: []byte(pass),
	}
}
