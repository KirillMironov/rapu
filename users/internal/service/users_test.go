package service

import (
	"context"
	"github.com/KirillMironov/rapu/users/internal/domain"
	"github.com/KirillMironov/rapu/users/test/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testUsername = "Lisa"
	testEmail    = "lisa@gmail.com"
	testPassword = "qwerty"
	testToken    = "token"
)

var usersService = NewUsers(mock.UsersRepository{}, mock.JWT{})

func TestUsersService_SignUp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		username      string
		email         string
		password      string
		expectToken   bool
		expectedError error
	}{
		{testUsername, testEmail, testPassword, true, nil},
		{"", testEmail, testPassword, false, domain.ErrEmptyParameters},
		{testUsername, "", testPassword, false, domain.ErrEmptyParameters},
		{testUsername, testEmail, "", false, domain.ErrEmptyParameters},
		{"", "", testPassword, false, domain.ErrEmptyParameters},
		{"", testEmail, "", false, domain.ErrEmptyParameters},
		{testUsername, "", "", false, domain.ErrEmptyParameters},
		{"", "", "", false, domain.ErrEmptyParameters},
	}

	for _, tc := range testCases {
		tc := tc

		token, err := usersService.SignUp(context.Background(), domain.User{
			Username: tc.username,
			Email:    tc.email,
			Password: tc.password,
		})
		assert.Equal(t, tc.expectedError, err)

		if tc.expectToken {
			assert.NotEmpty(t, token)
		} else {
			assert.Empty(t, token)
		}
	}
}

func TestUsersService_SignIn(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		email         string
		password      string
		expectToken   bool
		expectedError error
	}{
		{testEmail, testPassword, true, nil},
		{"", testPassword, false, domain.ErrEmptyParameters},
		{testEmail, "", false, domain.ErrEmptyParameters},
		{"", "", false, domain.ErrEmptyParameters},
	}

	for _, tc := range testCases {
		tc := tc

		token, err := usersService.SignIn(context.Background(), domain.User{
			Email:    tc.email,
			Password: tc.password,
		})
		assert.Equal(t, tc.expectedError, err)

		if tc.expectToken {
			assert.NotEmpty(t, token)
		} else {
			assert.Empty(t, token)
		}
	}
}

func TestUsersService_Authenticate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		token         string
		expectUserId  bool
		expectedError error
	}{
		{testToken, true, nil},
		{"", false, domain.ErrEmptyParameters},
	}

	for _, tc := range testCases {
		tc := tc

		userId, err := usersService.Authenticate(tc.token)
		assert.Equal(t, tc.expectedError, err)

		if tc.expectUserId {
			assert.NotEmpty(t, userId)
		} else {
			assert.Empty(t, userId)
		}
	}
}

func TestUsersService_UserExists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		userId            string
		expectedExistence bool
		expectedError     error
	}{
		{"1", true, nil},
		{"", false, domain.ErrEmptyParameters},
	}

	for _, tc := range testCases {
		tc := tc

		exists, err := usersService.UserExists(context.Background(), tc.userId)
		assert.Equal(t, tc.expectedError, err)
		assert.Equal(t, tc.expectedExistence, exists)
	}
}
