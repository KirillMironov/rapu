package service

import (
	"github.com/KirillMironov/rapu/domain"
	"github.com/KirillMironov/rapu/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

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

func (u *UsersService) SignIn(user domain.User) (string, error) {
	userId, hashedPassword, err := u.repository.GetByEmail(user)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return "", err
	}

	return u.tokenManager.GenerateAuthToken(userId)
}
