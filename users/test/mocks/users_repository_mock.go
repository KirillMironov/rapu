package mocks

import "github.com/KirillMironov/rapu/users/domain"

type UsersRepositoryMock struct{}

func (UsersRepositoryMock) Create(domain.User) (string, error) {
	return "", nil
}

func (UsersRepositoryMock) GetByEmail(string) (domain.User, error) {
	return domain.User{Password: "$2a$12$PRpg66gcvkLijyTzJtHVIeTucD/FAsvf/M8TWEt0O8GoJOdAkmRXK"}, nil
}
