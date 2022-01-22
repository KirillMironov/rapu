package service

import (
	"errors"
	"github.com/KirillMironov/rapu/users/domain"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repository   domain.UsersRepository
	tokenManager TokenManager
}

type TokenManager interface {
	Generate(userId string) (string, error)
	Verify(token string) (string, error)
}

func NewUsersService(repository domain.UsersRepository, tokenManager TokenManager) *UsersService {
	return &UsersService{
		repository:   repository,
		tokenManager: tokenManager,
	}
}

func (u *UsersService) SignUp(user domain.User) (string, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return "", domain.ErrEmptyParameters
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

	return u.tokenManager.Generate(userId)
}

func (u *UsersService) SignIn(input domain.User) (string, error) {
	if input.Email == "" || input.Password == "" {
		return "", domain.ErrEmptyParameters
	}

	user, err := u.repository.GetByEmail(input.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", domain.ErrWrongPassword
		}
		return "", err
	}

	return u.tokenManager.Generate(user.Id)
}

func (u *UsersService) Authenticate(token string) (string, error) {
	if token == "" {
		return "", domain.ErrEmptyParameters
	}

	return u.tokenManager.Verify(token)
}

func (u *UsersService) UserExists(userId string) (bool, error) {
	if userId == "" {
		return false, domain.ErrEmptyParameters
	}

	return u.repository.CheckExistence(userId)
}
