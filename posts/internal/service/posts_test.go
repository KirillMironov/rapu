package service

import (
	"github.com/KirillMironov/rapu/posts/domain"
	"github.com/KirillMironov/rapu/posts/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	userId  = "1"
	message = "Hello"
)

var svc = NewPostsService(&mocks.PostsRepositoryMock{}, 10)

func TestPostsService_Create(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		post          domain.Post
		expectedError error
	}{
		{domain.Post{UserId: userId, Message: message}, nil},
		{domain.Post{UserId: "", Message: message}, errNotEnoughArgs},
		{domain.Post{UserId: userId, Message: ""}, errNotEnoughArgs},
		{domain.Post{}, errNotEnoughArgs},
	}

	for _, tc := range testCases {
		tc := tc
		err := svc.Create(tc.post)
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
		{"", errNotEnoughArgs},
	}

	for _, tc := range testCases {
		tc := tc
		_, err := svc.GetByUserId(tc.userId, "", 0)
		assert.Equal(t, tc.expectedError, err)
	}
}
