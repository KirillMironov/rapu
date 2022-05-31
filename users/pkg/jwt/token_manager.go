package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type TokenManager struct {
	JWTKey   string
	TokenTTL time.Duration
}

func NewTokenManager(JWTKey string, tokenTTL time.Duration) (*TokenManager, error) {
	if JWTKey == "" {
		return nil, errors.New("JWT key was not provided")
	}
	return &TokenManager{JWTKey: JWTKey, TokenTTL: tokenTTL}, nil
}

func (m TokenManager) Generate(userId string) (string, error) {
	currentTime := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(m.TokenTTL)),
		IssuedAt:  jwt.NewNumericDate(currentTime),
	})

	return token.SignedString([]byte(m.JWTKey))
}

func (m TokenManager) Verify(token string) (string, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.JWTKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims["sub"].(string), nil
	}
	return "", errors.New("token is not valid")
}
