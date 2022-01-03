package mocks

import (
	"github.com/KirillMironov/rapu/users/domain"
)

type UsersServiceMock struct{}

func (UsersServiceMock) SignUp(user domain.User) (string, error) {
	return "token", nil
}

func (UsersServiceMock) SignIn(user domain.User) (string, error) {
	return "token", nil
}

func (UsersServiceMock) Authenticate(token string) (string, error) {
	return "token", nil
}
