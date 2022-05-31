package service

import (
	"errors"
	"github.com/KirillMironov/rapu/users/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	repository   UsersRepository
	tokenManager TokenManager
}

type UsersRepository interface {
	Create(domain.User) (string, error)
	GetByEmail(email string) (domain.User, error)
	CheckExistence(userId string) (bool, error)
}

type TokenManager interface {
	Generate(userId string) (string, error)
	Verify(token string) (string, error)
}

func NewUsers(repository UsersRepository, tokenManager TokenManager) *Users {
	return &Users{
		repository:   repository,
		tokenManager: tokenManager,
	}
}

func (u *Users) SignUp(user domain.User) (string, error) {
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

func (u *Users) SignIn(input domain.User) (string, error) {
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

func (u *Users) Authenticate(token string) (string, error) {
	if token == "" {
		return "", domain.ErrEmptyParameters
	}

	return u.tokenManager.Verify(token)
}

func (u *Users) UserExists(userId string) (bool, error) {
	if userId == "" {
		return false, domain.ErrEmptyParameters
	}

	return u.repository.CheckExistence(userId)
}
