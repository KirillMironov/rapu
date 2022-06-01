package mock

import "github.com/KirillMironov/rapu/users/internal/domain"

type UsersService struct{}

func (UsersService) SignUp(domain.User) (string, error) {
	return "token", nil
}

func (UsersService) SignIn(domain.User) (string, error) {
	return "token", nil
}

func (UsersService) Authenticate(string) (string, error) {
	return "token", nil
}
