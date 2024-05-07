package jwt

import (
	"encoding/base64"
	"errors"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var ErrTokenInvalid = errors.New("token invalid")

type Tokenizer struct {
	method jwt.SigningMethod
	secret []byte
}

func NewTokenizer(secret []byte) *Tokenizer {
	return &Tokenizer{method: jwt.SigningMethodHS256, secret: secret}
}

func (t *Tokenizer) MakeToken(user model.UserCredentials, duration time.Duration) (string, error) {
	token := jwt.New(t.method)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = base64.StdEncoding.EncodeToString(user.Id)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenStr, err := token.SignedString(t.secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (t *Tokenizer) Verify(token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return t.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		return claims, nil
	} else {
		return nil, ErrTokenInvalid
	}
}
