package service

import (
	"errors"
	"github.com/KirillMironov/rapu/users/domain"
	"github.com/KirillMironov/rapu/users/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

var errNotEnoughArgs = errors.New("not enough arguments")

type UsersService struct {
	repository   domain.UsersRepository
	tokenManager auth.TokenManager
}

func NewUsersService(repository domain.UsersRepository, tokenManager auth.TokenManager) *UsersService {
	return &UsersService{
		repository:   repository,
		tokenManager: tokenManager,
	}
}

func (u *UsersService) SignUp(user domain.User) (string, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return "", errNotEnoughArgs
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hash)

	userId, err := u.repository.Create(user)
	if err != nil {
		return "", err
	}

	return u.tokenManager.GenerateAuthToken(userId)
}

func (u *UsersService) SignIn(input domain.User) (string, error) {
	if input.Email == "" || input.Password == "" {
		return "", errNotEnoughArgs
	}

	user, err := u.repository.GetByEmail(input.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", err
	}

	return u.tokenManager.GenerateAuthToken(user.Id)
}

func (u *UsersService) Authenticate(token string) (string, error) {
	if token == "" {
		return "", errNotEnoughArgs
	}

	return u.tokenManager.VerifyAuthToken(token)
}
