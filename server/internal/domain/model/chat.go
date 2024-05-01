package model

type ChatType int32

const (
	CTUnknown ChatType = iota
	CTPersonal
)

type Chat struct {
	Id      []byte       `bson:"_id"`
	Type    ChatType     `bson:"type"`
	Members []ChatMember `bson:"members"`
}

type ChatMember struct {
	Uid []byte `bson:"uid"`
}

type UserChats struct {
	Uid   []byte   `bson:"_id"`
	Chats [][]byte `bson:"chats"`
}

type ChatMessage struct {
	Id        []byte
	Cid       []byte
	Uid       []byte
	Text      string
	Timestamp int64
}
