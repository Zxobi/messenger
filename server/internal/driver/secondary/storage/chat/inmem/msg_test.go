package inmem

import (
	"bytes"
	"context"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"slices"
	"strconv"
	"testing"
	"time"
)

func TestMessageStorage_Save(t *testing.T) {
	storage := NewMessageStorage()

	for i := 1; i <= 10; i++ {
		cid := [16]byte(uuid.New())
		uid := [16]byte(uuid.New())
		text := "test-text" + strconv.Itoa(i)

		for j := 1; j < 5; j++ {
			actualR, err := storage.Save(context.Background(), cid[:], uid[:], text)
			require.NoError(t, err, "save should not error")

			assert.NotEmpty(t, actualR.Id, "returned id not empty")
			assert.WithinDuration(t, time.Now(), time.UnixMilli(actualR.Timestamp), time.Second,
				"returned timestamp match expected")

			assert.ElementsMatch(t, cid, actualR.Cid, "returned cid match expected")
			assert.ElementsMatch(t, uid, actualR.Uid, "returned uid match expected")
			assert.Equal(t, text, actualR.Text, "returned text match expected")

			assert.Len(t, storage.hash, i, "hash len correct")

			chatMsgs, ok := storage.hash[cid]
			require.True(t, ok, "hash contains chat")
			require.Len(t, chatMsgs, j, "saved chat messages len correct")

			idx := slices.IndexFunc(chatMsgs, func(msg model.ChatMessage) bool {
				return bytes.Equal(msg.Id, actualR.Id)
			})
			require.NotEqual(t, -1, idx, "returned message should be in chat msgs")
			actualS := chatMsgs[idx]
			assert.WithinDuration(t, time.Now(), time.UnixMilli(actualS.Timestamp), time.Second,
				"saved timestamp match expected")

			assert.ElementsMatch(t, cid, actualS.Cid, "saved cid match expected")
			assert.ElementsMatch(t, uid, actualS.Uid, "saved uid match expected")
			assert.Equal(t, text, actualS.Text, "saved text match expected")
		}
	}
}
