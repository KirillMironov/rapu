package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	userId   = "111"
	key      = "abc123"
	tokenTTL = time.Minute * 60
)

func TestJWT(t *testing.T) {
	t.Parallel()

	_, err := NewJWT("", tokenTTL)
	assert.Error(t, err)

	_, err = NewJWT(key, tokenTTL)
	assert.NoError(t, err)
}

func TestJWT_Generate(t *testing.T) {
	t.Parallel()

	jwtService, err := NewJWT(key, tokenTTL)
	assert.NoError(t, err)

	token, err := jwtService.Generate(userId)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	token, err = jwtService.Generate("")
	assert.NotEmpty(t, token)
	assert.NoError(t, err)
}

func TestJWT_Verify(t *testing.T) {
	t.Parallel()

	jwtService, err := NewJWT(key, tokenTTL)
	assert.NoError(t, err)

	// token TTL = 60 Min
	token, err := jwtService.Generate(userId)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	id, err := jwtService.Verify(token)
	assert.Equal(t, userId, id)
	assert.NoError(t, err)

	// token TTL = -60 Min
	jwtService, err = NewJWT(key, -60*time.Minute)
	assert.NoError(t, err)

	token, err = jwtService.Generate(userId)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	id, err = jwtService.Verify(token)
	assert.Empty(t, id)
	assert.Error(t, err)
}
