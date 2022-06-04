package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWTManager struct {
	JWTKey   string
	TokenTTL time.Duration
}

func NewJWTManager(JWTKey string, tokenTTL time.Duration) (*JWTManager, error) {
	if JWTKey == "" {
		return nil, errors.New("JWT key was not provided")
	}
	return &JWTManager{JWTKey: JWTKey, TokenTTL: tokenTTL}, nil
}

func (jm JWTManager) Generate(userId string) (string, error) {
	currentTime := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(jm.TokenTTL)),
		IssuedAt:  jwt.NewNumericDate(currentTime),
	})

	return token.SignedString([]byte(jm.JWTKey))
}

func (jm JWTManager) Verify(token string) (string, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jm.JWTKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims["sub"].(string), nil
	}
	return "", errors.New("token is not valid")
}
