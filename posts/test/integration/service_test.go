//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"github.com/KirillMironov/rapu/posts/domain"
	"github.com/KirillMironov/rapu/posts/internal/delivery/proto"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

const (
	userId  = "1"
	message = "Hello"
)

var ctx = context.Background()

func TestPosts_Create(t *testing.T) {
	client, closeConn := newClient(t)
	defer closeConn()

	_, err := client.Create(ctx, &proto.CreateRequest{
		UserId:  userId,
		Message: message,
	})
	assert.NoError(t, err)

	_, err = client.Create(ctx, &proto.CreateRequest{
		UserId:  userId,
		Message: message,
	})
	assert.NoError(t, err)

	_, err = client.Create(ctx, &proto.CreateRequest{ // another UserId
		UserId:  "2",
		Message: message,
	})
	assert.NoError(t, err)

	resp, err := client.GetByUserId(ctx, &proto.GetByUserIdRequest{
		UserId: userId,
	})
	assert.NoError(t, err)

	var posts []domain.Post
	err = json.Unmarshal(resp.GetPosts(), &posts)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)

	resp, err = client.GetByUserId(ctx, &proto.GetByUserIdRequest{ // another UserId
		UserId: "2",
	})
	assert.NoError(t, err)

	posts = []domain.Post{}
	err = json.Unmarshal(resp.GetPosts(), &posts)
	assert.NoError(t, err)
	assert.Len(t, posts, 1)
}

func TestPosts_GetByUserId(t *testing.T) {
	client, closeConn := newClient(t)
	defer closeConn()

	for i := 0; i < maxLimit; i++ {
		_, err := client.Create(ctx, &proto.CreateRequest{
			UserId:  userId,
			Message: message,
		})
		assert.NoError(t, err)
	}

	resp, err := client.GetByUserId(ctx, &proto.GetByUserIdRequest{
		UserId: userId,
	})
	assert.NoError(t, err)

	var posts []domain.Post
	err = json.Unmarshal(resp.GetPosts(), &posts)
	assert.NoError(t, err)

	// descending order
	for i := 0; i < len(posts)-1; i++ {
		assert.True(t, posts[i].CreatedAt.After(posts[i+1].CreatedAt) || posts[i].CreatedAt.Equal(posts[i+1].CreatedAt))
	}
}

func TestPosts_GetByUserId_pagination_offset(t *testing.T) {
	client, closeConn := newClient(t)
	defer closeConn()

	for i := 0; i < maxLimit; i++ {
		_, err := client.Create(ctx, &proto.CreateRequest{
			UserId:  userId,
			Message: strconv.Itoa(i),
		})
		assert.NoError(t, err)
	}

	// default offset
	resp, err := client.GetByUserId(ctx, &proto.GetByUserIdRequest{
		UserId: userId,
	})
	assert.NoError(t, err)

	var posts []domain.Post
	err = json.Unmarshal(resp.GetPosts(), &posts)
	assert.NoError(t, err)
	assert.Len(t, posts, maxLimit)

	// manual offset
	resp, err = client.GetByUserId(ctx, &proto.GetByUserIdRequest{
		UserId: userId,
		Offset: posts[0].Id.Hex(),
	})
	assert.NoError(t, err)

	var postsWithOffset []domain.Post
	err = json.Unmarshal(resp.GetPosts(), &postsWithOffset)
	assert.NoError(t, err)
	assert.Len(t, postsWithOffset, 9)

	assert.Equal(t, posts[1], postsWithOffset[0])
}

func TestPosts_GetByUserId_pagination_limit(t *testing.T) {
	client, closeConn := newClient(t)
	defer closeConn()

	for i := 0; i < maxLimit*2; i++ {
		_, err := client.Create(ctx, &proto.CreateRequest{
			UserId:  userId,
			Message: strconv.Itoa(i),
		})
		assert.NoError(t, err)
	}

	testCases := []struct {
		limit          int64
		expectedLength int
	}{
		{0, maxLimit},
		{-1, maxLimit},
		{maxLimit * 2, maxLimit},
		{maxLimit / 2, maxLimit / 2},
	}

	for _, tc := range testCases {
		resp, err := client.GetByUserId(ctx, &proto.GetByUserIdRequest{
			UserId: userId,
			Limit:  tc.limit,
		})
		assert.NoError(t, err)

		var posts []domain.Post
		err = json.Unmarshal(resp.GetPosts(), &posts)
		assert.NoError(t, err)
		assert.Len(t, posts, tc.expectedLength)
	}
}
