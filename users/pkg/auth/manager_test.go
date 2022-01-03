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

	_, err := NewManager("", tokenTTL)
	assert.Error(t, err)

	_, err = NewManager(JWTKey, tokenTTL)
	assert.NoError(t, err)
}

func TestManager_GenerateAuthToken(t *testing.T) {
	t.Parallel()

	tokenManager, err := NewManager(JWTKey, tokenTTL)
	assert.NoError(t, err)

	token, err := tokenManager.GenerateAuthToken(userId)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	token, err = tokenManager.GenerateAuthToken("")
	assert.NotEmpty(t, token)
	assert.NoError(t, err)
}

func TestManager_VerifyAuthToken(t *testing.T) {
	t.Parallel()

	tokenManager, err := NewManager(JWTKey, tokenTTL)
	assert.NoError(t, err)

	// token TTL = 60 Min
	token, err := tokenManager.GenerateAuthToken(userId)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	id, err := tokenManager.VerifyAuthToken(token)
	assert.Equal(t, userId, id)
	assert.NoError(t, err)

	// token TTL = -60 Min
	tokenManager, err = NewManager(JWTKey, time.Duration(-60*int64(time.Minute)))
	assert.NoError(t, err)

	token, err = tokenManager.GenerateAuthToken(userId)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	id, err = tokenManager.VerifyAuthToken(token)
	assert.Empty(t, id)
	assert.Error(t, err)
}
