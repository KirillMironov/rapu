package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	userId   = "111"
	JWTKey   = "abc123"
	tokenTTL = time.Minute * 60
)

func TestManager(t *testing.T) {
	t.Parallel()

	_, err := NewTokenManager("", tokenTTL)
	assert.Error(t, err)

	_, err = NewTokenManager(JWTKey, tokenTTL)
	assert.NoError(t, err)
}

func TestManager_GenerateAuthToken(t *testing.T) {
	t.Parallel()

	tokenManager, err := NewTokenManager(JWTKey, tokenTTL)
	assert.NoError(t, err)

	token, err := tokenManager.Generate(userId)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	token, err = tokenManager.Generate("")
	assert.NotEmpty(t, token)
	assert.NoError(t, err)
}

func TestManager_VerifyAuthToken(t *testing.T) {
	t.Parallel()

	tokenManager, err := NewTokenManager(JWTKey, tokenTTL)
	assert.NoError(t, err)

	// token TTL = 60 Min
	token, err := tokenManager.Generate(userId)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	id, err := tokenManager.Verify(token)
	assert.Equal(t, userId, id)
	assert.NoError(t, err)

	// token TTL = -60 Min
	tokenManager, err = NewTokenManager(JWTKey, time.Duration(-60*int64(time.Minute)))
	assert.NoError(t, err)

	token, err = tokenManager.Generate(userId)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	id, err = tokenManager.Verify(token)
	assert.Empty(t, id)
	assert.Error(t, err)
}
