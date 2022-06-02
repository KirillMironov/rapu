package mock

import (
	"context"
	"github.com/KirillMironov/rapu/users/internal/domain"
)

type UsersRepository struct{}

func (UsersRepository) Create(context.Context, domain.User) (string, error) {
	return "", nil
}

func (UsersRepository) GetByEmail(context.Context, string) (domain.User, error) {
	return domain.User{Password: "$2a$12$PRpg66gcvkLijyTzJtHVIeTucD/FAsvf/M8TWEt0O8GoJOdAkmRXK"}, nil
}

func (UsersRepository) CheckExistence(context.Context, string) (bool, error) {
	return true, nil
}
