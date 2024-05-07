package converter

import (
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/pkg/cutils"
	protocolv1 "github.com/dvid-messanger/protos/gen/protocol"
)

func ChatToDTO(c *model.Chat) *protocolv1.Chat {
	proto := protocolv1.Chat{
		Id:          c.Id,
		Type:        protocolv1.ChatType(c.Type),
		ChatMembers: make([]*protocolv1.ChatMember, 0, len(c.Members)),
	}

	for _, member := range c.Members {
		proto.ChatMembers = append(proto.GetChatMembers(), &protocolv1.ChatMember{
			Uid: member.Uid,
		})
	}
	return &proto
}

func ChatFromDTO(c *protocolv1.Chat) *model.Chat {
	chat := model.Chat{
		Id:      c.GetId(),
		Type:    model.ChatType(c.GetType()),
		Members: make([]model.ChatMember, 0, len(c.GetChatMembers())),
	}

	for _, member := range c.GetChatMembers() {
		chat.Members = append(chat.Members, model.ChatMember{
			Uid: member.GetUid(),
		})
	}

	return &chat
}

func ChatsToDTO(chats []model.Chat) []*protocolv1.Chat {
	chatsDTO := make([]*protocolv1.Chat, 0, len(chats))
	for _, v := range chats {
		chatsDTO = append(chatsDTO, ChatToDTO(&v))
	}

	return chatsDTO
}

func ChatsFromDTO(chatsDTO []*protocolv1.Chat) []model.Chat {
	chats := make([]model.Chat, 0, len(chatsDTO))
	for _, v := range chatsDTO {
		chats = append(chats, *ChatFromDTO(v))
	}

	return chats
}

func ChatMessageToDTO(c *model.ChatMessage) *protocolv1.ChatMessage {
	return &protocolv1.ChatMessage{
		Id:        c.Id,
		Cid:       c.Cid,
		Uid:       c.Uid,
		Text:      c.Text,
		Timestamp: c.Timestamp,
	}
}

func ChatMessageFromDTO(c *protocolv1.ChatMessage) *model.ChatMessage {
	return &model.ChatMessage{
		Id:        c.GetId(),
		Cid:       c.GetCid(),
		Uid:       c.GetUid(),
		Text:      c.GetText(),
		Timestamp: c.GetTimestamp(),
	}
}

func ChatMessagesToDTO(c []model.ChatMessage) []*protocolv1.ChatMessage {
	return cutils.Map(c, func(msg model.ChatMessage) *protocolv1.ChatMessage {
		return ChatMessageToDTO(&msg)
	})
}

func ChatMessagesFromDTO(c []*protocolv1.ChatMessage) []model.ChatMessage {
	return cutils.Map(c, func(msg *protocolv1.ChatMessage) model.ChatMessage {
		return *ChatMessageFromDTO(msg)
	})
}
