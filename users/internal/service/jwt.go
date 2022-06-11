package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWT struct {
	key      []byte
	tokenTTL time.Duration
}

func NewJWT(key string, tokenTTL time.Duration) (*JWT, error) {
	if key == "" {
		return nil, errors.New("jwt key was not provided")
	}
	return &JWT{key: []byte(key), tokenTTL: tokenTTL}, nil
}

func (j JWT) Generate(userId string) (string, error) {
	var currentTime = time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(j.tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(currentTime),
	})

	return token.SignedString(j.key)
}

func (j JWT) Verify(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected token signing method: %v", token.Header["alg"])
		}
		return j.key, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims["sub"].(string), nil
	}
	return "", errors.New("token is not valid")
}
