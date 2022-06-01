package mock

import "github.com/KirillMironov/rapu/users/internal/domain"

type UsersServiceMock struct{}

func (UsersServiceMock) SignUp(domain.User) (string, error) {
	return "token", nil
}

func (UsersServiceMock) SignIn(domain.User) (string, error) {
	return "token", nil
}

func (UsersServiceMock) Authenticate(string) (string, error) {
	return "token", nil
}
