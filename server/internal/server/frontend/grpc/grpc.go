package grpc

import (
	"context"
	"github.com/dvid-messanger/internal/domain/converter"
	"github.com/dvid-messanger/internal/domain/model"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Notifier interface {
	NewMessage(ctx context.Context, message *model.ChatMessage) error
	NewChat(ctx context.Context, message *model.Chat) error
}

type serverApi struct {
	frontendv1.UnimplementedNotifierServer
	notifier Notifier
}

func Register(gRpc *grpc.Server, notifier Notifier) {
	frontendv1.RegisterNotifierServer(gRpc, &serverApi{notifier: notifier})
}

func (s *serverApi) NewMessage(ctx context.Context, req *frontendv1.NewMessageRequest) (*frontendv1.NewMessageResponse, error) {
	if err := validateNewMessage(req); err != nil {
		return nil, err
	}

	err := s.notifier.NewMessage(ctx, converter.ChatMessageFromDTO(req.GetMessage()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &frontendv1.NewMessageResponse{}, nil
}

func (s *serverApi) NewChat(ctx context.Context, req *frontendv1.NewChatRequest) (*frontendv1.NewChatResponse, error) {
	if err := validateNewChat(req); err != nil {
		return nil, err
	}

	err := s.notifier.NewChat(ctx, converter.ChatFromDTO(req.GetChat()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &frontendv1.NewChatResponse{}, nil
}

func validateNewMessage(req *frontendv1.NewMessageRequest) error {
	_ = req
	return nil
}

func validateNewChat(req *frontendv1.NewChatRequest) error {
	_ = req
	return nil
}
