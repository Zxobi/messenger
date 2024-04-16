package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/client/chat"
	"github.com/dvid-messanger/internal/domain/converter"
	"github.com/dvid-messanger/internal/lib/id"
	"github.com/dvid-messanger/internal/lib/proto"
	"github.com/dvid-messanger/internal/server/frontend/ws"
	frontendv1 "github.com/dvid-messanger/protos/gen/go/frontend"
	"log/slog"
)

func (h *Handler) handleEcho(c *ws.Client, echo *frontendv1.UpstreamEcho) error {
	const op = "handler.handleEcho"
	const dt = frontendv1.DownstreamType_D_ECHO

	downstream := frontendv1.DownstreamEcho{Content: echo.Content}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}

func (h *Handler) handleGetUser(c *ws.Client, upstream *frontendv1.UpstreamGetUser) error {
	const op = "handler.handleGetUser"
	const dt = frontendv1.DownstreamType_D_GET_USER

	if _, err := h.requireLogin(c); err != nil {
		if errors.Is(err, ErrNoAuth) {
			c.Send(mustMakeUnauthorized(dt))
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	user, err := h.user.User(context.TODO(), upstream.GetUid())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamGetUser{User: converter.UserToDTO(user)}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}

func (h *Handler) handleCurUser(c *ws.Client, _ *frontendv1.UpstreamCurUser) error {
	const op = "handler.handleGetUser"
	const dt = frontendv1.DownstreamType_D_CUR_USER

	uid, err := h.requireLogin(c)
	if err != nil {
		if errors.Is(err, ErrNoAuth) {
			c.Send(mustMakeUnauthorized(dt))
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	user, err := h.user.User(context.TODO(), uid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamCurUser{User: converter.UserToDTO(user)}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}

func (h *Handler) handleGetUsers(c *ws.Client, _ *frontendv1.UpstreamGetUsers) error {
	const op = "handler.handleGetUser"
	const dt = frontendv1.DownstreamType_D_GET_USERS

	if _, err := h.requireLogin(c); err != nil {
		if errors.Is(err, ErrNoAuth) {
			c.Send(mustMakeUnauthorized(dt))
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	user, err := h.user.Users(context.TODO())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamGetUsers{Users: converter.UsersToDTO(user)}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}

func (h *Handler) handleRegUser(c *ws.Client, upstream *frontendv1.UpstreamRegUser) error {
	const op = "handler.handleRegUser"
	const dt = frontendv1.DownstreamType_D_REG_USER

	user, err := h.user.Create(context.TODO(), upstream.GetEmail(), upstream.GetBio())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = h.auth.Register(context.TODO(), user.Id, upstream.GetEmail(), upstream.GetPassword()); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamRegUser{User: converter.UserToDTO(user)}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}

func (h *Handler) handleLogin(c *ws.Client, upstream *frontendv1.UpstreamLogin) error {
	const op = "handler.handleLogin"
	log := h.log.With(slog.String("op", op))

	log.Debug("start client login")

	token, err := h.auth.Login(context.TODO(), upstream.GetEmail(), upstream.GetPassword())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	c.SetValue(authKey, token)

	uid, err := h.requireLogin(c)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	usr, err := h.user.User(context.TODO(), uid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	chats, err := h.chat.UserChats(context.TODO(), uid)
	if err != nil && !errors.Is(err, chat.ErrUserChatsNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = h.registry.Register(c, usr, chats)
	if err != nil {
		log.Error("failed to register client")
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamLogin{Token: token}
	msg, err := proto.MarshalDownstream(&downstream, frontendv1.DownstreamType_D_LOGIN, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("client logged in " + id.String(usr.Id))
	c.Send(msg)
	return nil
}

func (h *Handler) handleGetChat(c *ws.Client, upstream *frontendv1.UpstreamGetChat) error {
	const op = "handler.handleChat"
	const dt = frontendv1.DownstreamType_D_GET_CHAT

	if _, err := h.requireLogin(c); err != nil {
		if errors.Is(err, ErrNoAuth) {
			c.Send(mustMakeUnauthorized(dt))
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	rChat, err := h.chat.Chat(context.TODO(), upstream.GetCid())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamGetChat{Chat: converter.ChatToDTO(rChat)}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}

func (h *Handler) handleGetUserChats(c *ws.Client, _ *frontendv1.UpstreamGetUserChats) error {
	const op = "handler.handleChat"
	const dt = frontendv1.DownstreamType_D_GET_USER_CHATS

	uid, err := h.requireLogin(c)
	if err != nil {
		if errors.Is(err, ErrNoAuth) {
			c.Send(mustMakeUnauthorized(dt))
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	rChat, err := h.chat.UserChats(context.TODO(), uid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamGetUserChats{Chats: converter.ChatsToDTO(rChat)}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}

func (h *Handler) handleCreateChat(c *ws.Client, upstream *frontendv1.UpstreamCreateChat) error {
	const op = "handler.handleCreateChat"
	const dt = frontendv1.DownstreamType_D_CREATE_CHAT

	uid, err := h.requireLogin(c)
	if err != nil {
		if errors.Is(err, ErrNoAuth) {
			c.Send(mustMakeUnauthorized(dt))
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	rChat, err := h.chat.Create(context.TODO(), uid, upstream.GetUid())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamCreateChat{Chat: converter.ChatToDTO(rChat)}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}

func (h *Handler) handleSendMessage(c *ws.Client, upstream *frontendv1.UpstreamSendMessage) error {
	const op = "handler.handleSendMessage"
	const dt = frontendv1.DownstreamType_D_SEND_MESSAGE
	log := h.log.With(slog.String("op", op))

	uid, err := h.requireLogin(c)
	if err != nil {
		if errors.Is(err, ErrNoAuth) {
			c.Send(mustMakeUnauthorized(dt))
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	message, err := h.chat.SendMessage(context.TODO(), upstream.GetCid(), uid, upstream.GetText())
	if err != nil {
		log.Error("failed to send message")
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamSendMessage{Message: converter.ChatMessageToDTO(message)}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}

func (h *Handler) handleChatMessages(c *ws.Client, upstream *frontendv1.UpstreamChatMessages) error {
	const op = "handler.handleCreateChat"
	const dt = frontendv1.DownstreamType_D_CHAT_MESSAGES

	_, err := h.requireLogin(c)
	if err != nil {
		if errors.Is(err, ErrNoAuth) {
			c.Send(mustMakeUnauthorized(dt))
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	messages, err := h.chat.Messages(context.TODO(), upstream.GetCid())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream := frontendv1.DownstreamChatMessages{Messages: converter.ChatMessagesToDTO(messages)}
	msg, err := proto.MarshalDownstream(&downstream, dt, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.Send(msg)
	return nil
}
