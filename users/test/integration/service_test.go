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

	resp, err := client.SignUp(ctx, &proto.SignUpRequest{ // create user
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.GetAccessToken())

	resp, err = client.SignUp(ctx, &proto.SignUpRequest{ // empty credentials
		Username: "",
		Email:    "",
		Password: "",
	})
	assert.Error(t, err)
	assert.Empty(t, resp.GetAccessToken())
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	resp, err = client.SignUp(ctx, &proto.SignUpRequest{ // user already exists
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.Error(t, err)
	assert.Empty(t, resp.GetAccessToken())
	st, _ = status.FromError(err)
	assert.Equal(t, codes.AlreadyExists, st.Code())
}

func Test_SignIn(t *testing.T) {
	var (
		client = newClient(t)
		ctx    = context.Background()
	)

	resp, err := client.SignIn(ctx, &proto.SignInRequest{ // user doesn't exist
		Email:    testEmail,
		Password: testPassword,
	})
	assert.Error(t, err)
	assert.Empty(t, resp.GetAccessToken())

	resp, err = client.SignUp(ctx, &proto.SignUpRequest{ // create user
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.GetAccessToken())

	resp, err = client.SignIn(ctx, &proto.SignInRequest{ // sign in
		Email:    testEmail,
		Password: testPassword,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.GetAccessToken())

	resp, err = client.SignIn(ctx, &proto.SignInRequest{ // empty credentials
		Email:    "",
		Password: "",
	})
	assert.Error(t, err)
	assert.Empty(t, resp.GetAccessToken())
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	resp, err = client.SignIn(ctx, &proto.SignInRequest{ // invalid credentials
		Email:    "wrong email",
		Password: "wrong password",
	})
	assert.Error(t, err)
	assert.Empty(t, resp.GetAccessToken())
	st, _ = status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())

	resp, err = client.SignIn(ctx, &proto.SignInRequest{ // wrong password
		Email:    testEmail,
		Password: "wrong password",
	})
	assert.Error(t, err)
	assert.Empty(t, resp.GetAccessToken())
	st, _ = status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func Test_Authenticate(t *testing.T) {
	var (
		client = newClient(t)
		ctx    = context.Background()
	)

	resp, err := client.SignUp(ctx, &proto.SignUpRequest{ // create user
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.GetAccessToken())

	var token = resp.GetAccessToken()

	authResp, err := client.Authenticate(ctx, &proto.AuthRequest{AccessToken: token}) // get userId from token
	assert.NoError(t, err)
	assert.NotEmpty(t, authResp.GetUserId())

	authResp, err = client.Authenticate(ctx, &proto.AuthRequest{AccessToken: ""}) // empty token
	assert.Error(t, err)
	assert.Empty(t, authResp.GetUserId())
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	authResp, err = client.Authenticate(ctx, &proto.AuthRequest{AccessToken: "token"}) // invalid token
	assert.Error(t, err)
	assert.Empty(t, authResp.GetUserId())
	st, _ = status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func Test_UserExists(t *testing.T) {
	var (
		client = newClient(t)
		ctx    = context.Background()
	)

	resp, _ := client.SignUp(ctx, &proto.SignUpRequest{ // create user
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})

	authResp, _ := client.Authenticate(ctx, &proto.AuthRequest{AccessToken: resp.GetAccessToken()}) // get userId from token

	existsResp, err := client.UserExists(ctx, &proto.UserExistsRequest{UserId: authResp.GetUserId()})
	assert.NoError(t, err)
	assert.True(t, existsResp.GetExists())

	existsResp, err = client.UserExists(ctx, &proto.UserExistsRequest{UserId: ""}) // empty userId
	assert.Error(t, err)
	assert.False(t, existsResp.GetExists())
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	existsResp, err = client.UserExists(ctx, &proto.UserExistsRequest{UserId: "5"}) // user with userId "5" doesn't exist
	assert.Error(t, err)
	assert.False(t, existsResp.GetExists())
	st, _ = status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
}
