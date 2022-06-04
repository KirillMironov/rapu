package service

import (
	"context"
	"github.com/KirillMironov/rapu/posts/internal/domain"
	"github.com/KirillMironov/rapu/posts/test/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	userId  = "1"
	message = "Hello"
)

var (
	postsService = NewPosts(&mock.PostsRepository{}, 10, mock.Logger{})
	ctx          = context.Background()
)

func TestPostsService_Create(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		post          domain.Post
		expectedError error
	}{
		{domain.Post{UserId: userId, Message: message}, nil},
		{domain.Post{UserId: "", Message: message}, domain.ErrEmptyParameters},
		{domain.Post{UserId: userId, Message: ""}, domain.ErrEmptyParameters},
		{domain.Post{}, domain.ErrEmptyParameters},
	}

	for _, tc := range testCases {
		tc := tc
		err := postsService.Create(ctx, tc.post)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestPostsService_GetByUserId(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		userId        string
		expectedError error
	}{
		{userId, nil},
		{"", domain.ErrEmptyParameters},
	}

	for _, tc := range testCases {
		tc := tc
		_, err := postsService.GetByUserId(ctx, tc.userId, "", 0)
		assert.Equal(t, tc.expectedError, err)
	}
}
