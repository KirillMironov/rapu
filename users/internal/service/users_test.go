package service

import (
	"github.com/KirillMironov/rapu/users/domain"
	"github.com/KirillMironov/rapu/users/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testUsername = "Lisa"
	testEmail    = "lisa@gmail.com"
	testPassword = "qwerty"
	testToken    = "token"
)

var service = NewUsersService(mocks.UsersRepositoryMock{}, mocks.TokenManagerMock{})

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

		token, err := service.SignUp(domain.User{
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

		token, err := service.SignIn(domain.User{
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

		userId, err := service.Authenticate(tc.token)
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

		exists, err := service.UserExists(tc.userId)
		assert.Equal(t, tc.expectedError, err)
		assert.Equal(t, tc.expectedExistence, exists)
	}
}
