package service

import (
	"github.com/KirillMironov/rapu/users/domain"
	"github.com/KirillMironov/rapu/users/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	username = "Lisa"
	email    = "lisa@gmail.com"
	password = "qwerty"
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
		{username, email, password, true, nil},
		{"", email, password, false, errNotEnoughArgs},
		{username, "", password, false, errNotEnoughArgs},
		{username, email, "", false, errNotEnoughArgs},
		{"", "", password, false, errNotEnoughArgs},
		{"", email, "", false, errNotEnoughArgs},
		{username, "", "", false, errNotEnoughArgs},
		{"", "", "", false, errNotEnoughArgs},
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
		{email, password, true, nil},
		{"", password, false, errNotEnoughArgs},
		{email, "", false, errNotEnoughArgs},
		{"", "", false, errNotEnoughArgs},
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
		{"token", true, nil},
		{"", false, errNotEnoughArgs},
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
