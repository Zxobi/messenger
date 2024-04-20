package inmem

import (
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/id"
	"github.com/dvid-messanger/internal/storage/chat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestChatStorage_ChatSave(t *testing.T) {
	storage := NewChatStorage()

	for i := 1; i <= 10; i++ {
		uidFrom := [16]byte(uuid.New())
		uidTo := [16]byte(uuid.New())

		actualR, err := storage.Save(context.Background(), uidFrom[:], uidTo[:])
		require.NoError(t, err, "save should not error")

		assert.NotEmpty(t, actualR.Id, "returned id not empty")
		assert.Equal(t, model.CTPersonal, actualR.Type, "returned type match expected")

		require.Len(t, actualR.Members, 2, "returned chat have two members")

		assert.ElementsMatch(t, uidFrom, actualR.Members[0].Id, "returned first member id match expected")
		assert.ElementsMatch(t, uidTo, actualR.Members[1].Id, "returned second member id match expected")

		assert.Len(t, storage.hash, i, "hash len correct")

		actualS, ok := storage.hash[id.Id(actualR.Id)]
		assert.True(t, ok, "hash contains chat")
		assert.Equal(t, model.CTPersonal, actualS.Type, "saved type match expected")

		require.Len(t, actualS.Members, 2, "saved chat have two members")

		assert.ElementsMatch(t, uidFrom, actualS.Members[0].Id, "saved first member id match expected")
		assert.ElementsMatch(t, uidTo, actualS.Members[1].Id, "saved second member id match expected")
	}
}

func TestChatStorage_Save_ChatWithSameMembersErr(t *testing.T) {
	storage := NewChatStorage()

	uidFrom := [16]byte(uuid.New())
	uidTo := [16]byte(uuid.New())

	_, err := storage.Save(context.Background(), uidFrom[:], uidTo[:])
	require.NoError(t, err, "save should not error")

	_, err = storage.Save(context.Background(), uidFrom[:], uidTo[:])
	if assert.Error(t, err, "save with same members should error") {
		assert.ErrorIs(t, err, chat.ErrChatExists, "error type correct")
	}

	_, err = storage.Save(context.Background(), uidTo[:], uidFrom[:])
	if assert.Error(t, err, "save with same members should error") {
		assert.ErrorIs(t, err, chat.ErrChatExists, "error type correct")
	}
}

func TestChatStorage_Chat(t *testing.T) {
	storage := NewChatStorage()

	uid1 := [16]byte(uuid.New())
	uid2 := [16]byte(uuid.New())
	uid3 := [16]byte(uuid.New())
	saved1, err := storage.Save(context.Background(), uid1[:], uid2[:])
	require.NoError(t, err, "save should not error")
	saved2, err := storage.Save(context.Background(), uid2[:], uid3[:])
	require.NoError(t, err, "save should not error")

	actual1, err := storage.Chat(context.Background(), saved1.Id)
	require.NoError(t, err, "chat should not error")
	actual2, err := storage.Chat(context.Background(), saved2.Id)
	require.NoError(t, err, "chat should not error")

	assert.ElementsMatch(t, saved1.Id, actual1.Id, "id match saved")
	assert.ElementsMatch(t, saved1.Members, actual1.Members, "members match saved")

	assert.ElementsMatch(t, saved2.Id, actual2.Id, "id match saved")
	assert.ElementsMatch(t, saved2.Members, actual2.Members, "members match saved")
}

func TestChatStorage_UserChats(t *testing.T) {
	storage := NewChatStorage()

	uid1 := [16]byte(uuid.New())
	uid2 := [16]byte(uuid.New())
	uid3 := [16]byte(uuid.New())

	_, err := storage.Save(context.Background(), uid1[:], uid2[:])
	require.NoError(t, err, "save should not error")
	_, err = storage.Save(context.Background(), uid2[:], uid3[:])
	require.NoError(t, err, "save should not error")

	chats, err := storage.UserChats(context.Background(), uid1[:])
	require.Len(t, chats, 1, "user chats len correct")
	chats, err = storage.UserChats(context.Background(), uid2[:])
	require.Len(t, chats, 2, "user chats len correct")

	_, err = storage.Save(context.Background(), uid1[:], uid3[:])
	require.NoError(t, err, "save should not error")
	chats, err = storage.UserChats(context.Background(), uid1[:])
	require.Len(t, chats, 2, "user chats len correct")
}
