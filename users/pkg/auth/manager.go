package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenManager interface {
	GenerateAuthToken(userId string) (string, error)
	VerifyAuthToken(token string) (string, error)
}

type Manager struct {
	JWTKey   string
	TokenTTL time.Duration
}

func NewManager(JWTKey string, tokenTTL time.Duration) (*Manager, error) {
	if JWTKey == "" {
		return nil, errors.New("JWT key was not provided")
	}
	return &Manager{JWTKey: JWTKey, TokenTTL: tokenTTL}, nil
}

func (m Manager) GenerateAuthToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(m.TokenTTL).Unix(),
	})

	return token.SignedString([]byte(m.JWTKey))
}

func (m Manager) VerifyAuthToken(token string) (string, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.JWTKey), nil
	})
	if err != nil {
		return "", err
	}

	return tkn.Claims.(jwt.MapClaims)["sub"].(string), nil
}
