package model

type ChatType int32

const (
	CTUnknown ChatType = iota
	CTPersonal
)

type Chat struct {
	Id      []byte
	Type    ChatType
	Members []ChatMember
}

type ChatMember struct {
	Id []byte
}

type ChatMessage struct {
	Id        []byte
	Cid       []byte
	Uid       []byte
	Text      string
	Timestamp int64
}
