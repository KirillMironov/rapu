//go:build integration

package integration

import (
	"context"
	"github.com/KirillMironov/rapu/users/internal/delivery/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

const (
	testUsername = "Lisa"
	testEmail    = "lisa@gmail.com"
	testPassword = "qwerty"
)

func Test_SignUp(t *testing.T) {
	var (
		client = newClient(t)
		ctx    = context.Background()
	)

	testCases := []struct {
		name               string
		username           string
		email              string
		password           string
		expectAccessToken  bool
		expectError        bool
		expectedStatusCode codes.Code
	}{
		{
			name:               "create user",
			username:           testUsername,
			email:              testEmail,
			password:           testPassword,
			expectAccessToken:  true,
			expectError:        false,
			expectedStatusCode: codes.OK,
		},
		{
			name:               "empty credentials",
			username:           "",
			email:              "",
			password:           "",
			expectAccessToken:  false,
			expectError:        true,
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "user already exists",
			username:           testUsername,
			email:              testEmail,
			password:           testPassword,
			expectAccessToken:  false,
			expectError:        true,
			expectedStatusCode: codes.AlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.SignUp(ctx, &proto.SignUpRequest{
				Username: tc.username,
				Email:    tc.email,
				Password: tc.password,
			})
			assert.Equal(t, tc.expectError, err != nil)
			assert.Equal(t, tc.expectAccessToken, resp.GetAccessToken() != "")
			st, _ := status.FromError(err)
			assert.Equal(t, tc.expectedStatusCode, st.Code())
		})
	}
}

func Test_SignIn(t *testing.T) {
	var (
		client = newClient(t)
		ctx    = context.Background()
	)

	testCases := []struct {
		name               string
		email              string
		password           string
		expectAccessToken  bool
		expectError        bool
		expectedStatusCode codes.Code
		preparation        func(client proto.UsersClient) error
	}{
		{
			name:               "user doesn't exist",
			email:              testEmail,
			password:           testPassword,
			expectAccessToken:  false,
			expectError:        true,
			expectedStatusCode: codes.Unauthenticated,
		},
		{
			name:               "sign in",
			email:              testEmail,
			password:           testPassword,
			expectAccessToken:  true,
			expectError:        false,
			expectedStatusCode: codes.OK,
			preparation: func(client proto.UsersClient) error {
				_, err := client.SignUp(ctx, &proto.SignUpRequest{ // create user
					Username: testUsername,
					Email:    testEmail,
					Password: testPassword,
				})
				return err
			},
		},
		{
			name:               "empty credentials",
			email:              "",
			password:           "",
			expectAccessToken:  false,
			expectError:        true,
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "invalid credentials",
			email:              "q",
			password:           "q",
			expectAccessToken:  false,
			expectError:        true,
			expectedStatusCode: codes.Unauthenticated,
		},
		{
			name:               "invalid password",
			email:              testEmail,
			password:           "q",
			expectAccessToken:  false,
			expectError:        true,
			expectedStatusCode: codes.Unauthenticated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.preparation != nil {
				assert.NoError(t, tc.preparation(client))
			}
			resp, err := client.SignIn(ctx, &proto.SignInRequest{
				Email:    tc.email,
				Password: tc.password,
			})
			assert.Equal(t, tc.expectError, err != nil)
			assert.Equal(t, tc.expectAccessToken, resp.GetAccessToken() != "")
			st, _ := status.FromError(err)
			assert.Equal(t, tc.expectedStatusCode, st.Code())
		})
	}
}

func Test_Authenticate(t *testing.T) {
	var (
		client = newClient(t)
		ctx    = context.Background()
	)

	resp, err := client.SignUp(ctx, &proto.SignUpRequest{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.NoError(t, err)

	var token = resp.GetAccessToken()

	testCases := []struct {
		name               string
		accessToken        string
		expectUserId       bool
		expectError        bool
		expectedStatusCode codes.Code
	}{
		{
			name:               "userId from token",
			accessToken:        token,
			expectUserId:       true,
			expectError:        false,
			expectedStatusCode: codes.OK,
		},
		{
			name:               "empty token",
			accessToken:        "",
			expectUserId:       false,
			expectError:        true,
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "invalid token",
			accessToken:        "q",
			expectUserId:       false,
			expectError:        true,
			expectedStatusCode: codes.Unauthenticated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Authenticate(ctx, &proto.AuthRequest{
				AccessToken: tc.accessToken,
			})
			assert.Equal(t, tc.expectError, err != nil)
			assert.Equal(t, tc.expectUserId, resp.GetUserId() != "")
			st, _ := status.FromError(err)
			assert.Equal(t, tc.expectedStatusCode, st.Code())
		})
	}
}

func Test_UserExists(t *testing.T) {
	var (
		client = newClient(t)
		ctx    = context.Background()
	)

	resp, err := client.SignUp(ctx, &proto.SignUpRequest{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.NoError(t, err)
	authResp, err := client.Authenticate(ctx, &proto.AuthRequest{AccessToken: resp.GetAccessToken()})
	assert.NoError(t, err)

	var userId = authResp.GetUserId()

	testCases := []struct {
		name               string
		userId             string
		expectExistence    bool
		expectError        bool
		expectedStatusCode codes.Code
	}{
		{
			name:               "user exists",
			userId:             userId,
			expectExistence:    true,
			expectError:        false,
			expectedStatusCode: codes.OK,
		},
		{
			name:               "empty userId",
			userId:             "",
			expectExistence:    false,
			expectError:        true,
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "invalid userId",
			userId:             "-1",
			expectExistence:    false,
			expectError:        true,
			expectedStatusCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.UserExists(ctx, &proto.UserExistsRequest{
				UserId: tc.userId,
			})
			assert.Equal(t, tc.expectError, err != nil)
			assert.Equal(t, tc.expectExistence, resp.GetExists())
			st, _ := status.FromError(err)
			assert.Equal(t, tc.expectedStatusCode, st.Code())
		})
	}
}
