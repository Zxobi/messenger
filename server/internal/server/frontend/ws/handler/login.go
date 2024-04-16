package handler

import (
	"encoding/base64"
	"errors"
	"github.com/dvid-messanger/internal/server/frontend/ws"
)

const authKey = "auth"

var (
	ErrNoAuth   = errors.New("not authorized")
	ErrBadToken = errors.New("bad token")
)

func (h *Handler) requireLogin(c *ws.Client) ([]byte, error) {
	auth := c.GetValue(authKey)
	if auth == "" {
		return nil, ErrNoAuth
	}

	claims, err := h.tv.Verify(auth)
	if err != nil {
		return nil, err
	}

	uidStr, ok := claims["uid"]
	if !ok {
		return nil, ErrBadToken
	}

	uid, err := base64.StdEncoding.DecodeString(uidStr.(string))
	if err != nil {
		return nil, ErrBadToken
	}

	return uid, nil
}
