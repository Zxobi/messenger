package modelutil

import (
	"bytes"
	"github.com/dvid-messanger/internal/domain/model"
	"slices"
)

func HaveChatWith(userChats *model.UserChats, uid []byte) bool {
	return slices.ContainsFunc(userChats.Chats, func(userChat model.UserChat) bool {
		return bytes.Equal(userChat.Uid, uid)
	})
}

func AddPersonalChat(userChats *model.UserChats, cid []byte, uid []byte) *model.UserChat {
	userChat := &model.UserChat{Cid: cid, Type: model.CTPersonal, Uid: uid}
	userChats.Chats = append(userChats.Chats, *userChat)

	return userChat
}
