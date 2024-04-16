package chat

import (
	"context"
	"errors"
	"github.com/dvid-messanger/internal/domain/converter"
	"github.com/dvid-messanger/internal/domain/model"
	grpcutil "github.com/dvid-messanger/internal/lib/grpc"
	"github.com/dvid-messanger/internal/service/chat"
	chatv1 "github.com/dvid-messanger/protos/gen/go/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Chat interface {
	Create(ctx context.Context, fromUid []byte, toUid []byte) (*model.Chat, error)

	Chat(ctx context.Context, cid []byte) (*model.Chat, error)
	UserChats(ctx context.Context, uid []byte) ([]model.Chat, error)

	SendMessage(ctx context.Context, cid []byte, uid []byte, text string) (*model.ChatMessage, error)
	Messages(ctx context.Context, cid []byte) ([]model.ChatMessage, error)
}

type serverApi struct {
	chatv1.UnimplementedChatServiceServer
	chat Chat
}

func Register(gRpc *grpc.Server, chat Chat) {
	chatv1.RegisterChatServiceServer(gRpc, &serverApi{chat: chat})
}

func (s *serverApi) Create(ctx context.Context, req *chatv1.CreateChatRequest) (*chatv1.CreateChatResponse, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}

	c, err := s.chat.Create(ctx, req.GetFromUid(), req.GetToUid())
	if err != nil {
		if errors.Is(err, chat.ErrChatExists) {
			return nil, status.Error(codes.AlreadyExists, "chat already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &chatv1.CreateChatResponse{Chat: converter.ChatToDTO(c)}, nil
}

func (s *serverApi) Chat(ctx context.Context, req *chatv1.ChatRequest) (*chatv1.ChatResponse, error) {
	if err := validateChat(req); err != nil {
		return nil, err
	}

	c, err := s.chat.Chat(ctx, req.GetCid())
	if err != nil {
		if errors.Is(err, chat.ErrChatNotFound) {
			return nil, status.Error(codes.NotFound, "chat not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &chatv1.ChatResponse{Chat: converter.ChatToDTO(c)}, nil
}

func (s *serverApi) UserChats(ctx context.Context, req *chatv1.UserChatsRequest) (*chatv1.UserChatsResponse, error) {
	if err := validateUserChats(req); err != nil {
		return nil, err
	}

	chats, err := s.chat.UserChats(ctx, req.GetUid())
	if err != nil {
		if errors.Is(err, chat.ErrUserChatsNotFound) {
			return nil, status.Error(codes.NotFound, "user chats not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &chatv1.UserChatsResponse{Chats: converter.ChatsToDTO(chats)}, nil
}

func (s *serverApi) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*chatv1.SendMessageResponse, error) {
	if err := validateSendMessage(req); err != nil {
		return nil, err
	}

	msg, err := s.chat.SendMessage(ctx, req.GetCid(), req.GetUid(), req.GetText())
	if err != nil {
		if errors.Is(err, chat.ErrChatNotFound) {
			return nil, status.Error(codes.NotFound, "chat not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &chatv1.SendMessageResponse{Message: converter.ChatMessageToDTO(msg)}, nil
}
func (s *serverApi) Messages(ctx context.Context, req *chatv1.MessagesRequest) (*chatv1.MessagesResponse, error) {
	if err := validateMessages(req); err != nil {
		return nil, err
	}

	c, err := s.chat.Messages(ctx, req.GetCid())
	if err != nil {
		if errors.Is(err, chat.ErrChatNotFound) {
			return nil, status.Error(codes.NotFound, "chat not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &chatv1.MessagesResponse{Messages: converter.ChatMessagesToDTO(c)}, nil
}

func validateCreate(req *chatv1.CreateChatRequest) error {
	if err := grpcutil.ValidateId(req.GetToUid(), "toUid"); err != nil {
		return err
	}
	if err := grpcutil.ValidateId(req.GetFromUid(), "fromUid"); err != nil {
		return err
	}

	return nil
}

func validateChat(req *chatv1.ChatRequest) error {
	if err := grpcutil.ValidateId(req.GetCid(), "cid"); err != nil {
		return err
	}

	return nil
}

func validateUserChats(req *chatv1.UserChatsRequest) error {
	if err := grpcutil.ValidateId(req.GetUid(), "uid"); err != nil {
		return err
	}

	return nil
}

func validateSendMessage(req *chatv1.SendMessageRequest) error {
	if err := grpcutil.ValidateId(req.GetUid(), "uid"); err != nil {
		return err
	}
	if err := grpcutil.ValidateId(req.GetCid(), "cid"); err != nil {
		return err
	}
	if len(req.GetText()) == 0 {
		return status.Error(codes.InvalidArgument, "text is required")
	}

	return nil
}

func validateMessages(req *chatv1.MessagesRequest) error {
	if err := grpcutil.ValidateId(req.GetCid(), "cid"); err != nil {
		return err
	}

	return nil
}
