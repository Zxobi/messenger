package model

import "github.com/dvid-messanger/internal/pkg/cutils"

type ChatType int32

const (
	CTUnknown ChatType = iota
	CTPersonal
)

type Chat struct {
	Id      []byte       `bson:"_id"`
	Type    ChatType     `bson:"type"`
	Members []ChatMember `bson:"members,omitempty"`
}

func NewPersonalChat(cid []byte, uids ...[]byte) *Chat {
	return &Chat{
		Id:   cid,
		Type: CTPersonal,
		Members: cutils.Map(uids, func(uid []byte) ChatMember {
			return ChatMember{
				Uid: uid,
			}
		}),
	}
}

type ChatMember struct {
	Uid []byte `bson:"uid"`
}

type UserChats struct {
	Uid   []byte     `bson:"_id"`
	Chats []UserChat `bson:"chats,omitempty"`
}

type UserChat struct {
	Cid  []byte   `bson:"cid"`
	Type ChatType `bson:"type"`
	Uid  []byte   `bson:"uid"`
}

type ChatMessage struct {
	Id        []byte
	Cid       []byte
	Uid       []byte
	Text      string
	Timestamp int64
}
