//go:build integration

package integration

import (
	"context"
	"github.com/KirillMironov/rapu/posts/internal/delivery/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

const (
	userId  = "1"
	message = "Hello"
)

func TestPosts_Create(t *testing.T) {
	var (
		client = newClient(t)
		ctx    = context.Background()
	)

	testCases := []struct {
		name               string
		userId             string
		message            string
		expectError        bool
		expectedStatusCode codes.Code
	}{
		{
			name:               "create post",
			userId:             userId,
			message:            message,
			expectError:        false,
			expectedStatusCode: codes.OK,
		},
		{
			name:               "same post",
			userId:             userId,
			message:            message,
			expectError:        false,
			expectedStatusCode: codes.OK,
		},
		{
			name:               "empty parameters",
			userId:             "",
			message:            "",
			expectError:        true,
			expectedStatusCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.Create(ctx, &proto.CreateRequest{
				UserId:  tc.userId,
				Message: tc.message,
			})
			assert.Equal(t, tc.expectError, err != nil)
			st, _ := status.FromError(err)
			assert.Equal(t, tc.expectedStatusCode, st.Code())
		})
	}
}

func TestPosts_GetByUserId(t *testing.T) {
	var (
		client = newClient(t)
		ctx    = context.Background()
	)

	testCases := []struct {
		name               string
		userId             string
		message            string
		expectPosts        bool
		expectError        bool
		expectedStatusCode codes.Code
		preparation        func(client proto.PostsClient) error
	}{
		{
			name:               "no posts yet",
			userId:             userId,
			expectPosts:        false,
			expectError:        true,
			expectedStatusCode: codes.NotFound,
		},
		{
			name:               "get post",
			userId:             userId,
			expectPosts:        true,
			expectError:        false,
			expectedStatusCode: codes.OK,
			preparation: func(client proto.PostsClient) error {
				_, err := client.Create(ctx, &proto.CreateRequest{
					UserId:  userId,
					Message: message,
				})
				return err
			},
		},
		{
			name:               "empty parameters",
			userId:             "",
			expectPosts:        false,
			expectError:        true,
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "no posts with userId",
			userId:             "-1",
			expectPosts:        false,
			expectError:        true,
			expectedStatusCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.preparation != nil {
				assert.NoError(t, tc.preparation(client))
			}
			resp, err := client.GetByUserId(ctx, &proto.GetByUserIdRequest{
				UserId: tc.userId,
			})
			assert.Equal(t, tc.expectError, err != nil)
			st, _ := status.FromError(err)
			assert.Equal(t, tc.expectedStatusCode, st.Code())
			assert.Equal(t, tc.expectPosts, resp.GetPosts() != nil)
		})
	}
}
