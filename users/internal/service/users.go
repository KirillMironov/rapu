package service

import (
	"errors"
	"github.com/KirillMironov/rapu/users/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	repository   UsersRepository
	tokenManager TokenManager
	logger       Logger
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

type Logger interface {
	Error(args ...interface{})
}

func NewUsers(repository UsersRepository, tokenManager TokenManager, logger Logger) *Users {
	return &Users{
		repository:   repository,
		tokenManager: tokenManager,
		logger:       logger,
	}
}

func (u *Users) SignUp(user domain.User) (string, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return "", domain.ErrEmptyParameters
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(err)
		return "", err
	}
	user.Password = string(hash)

	userId, err := u.repository.Create(user)
	if err != nil {
		u.logger.Error(err)
		return "", err
	}

	token, err := u.tokenManager.Generate(userId)
	if err != nil {
		u.logger.Error(err)
		return "", err
	}

	return token, nil
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
		u.logger.Error(err)
		return "", err
	}

	token, err := u.tokenManager.Generate(user.Id)
	if err != nil {
		u.logger.Error(err)
		return "", err
	}

	return token, nil
}

func (u *Users) Authenticate(token string) (string, error) {
	if token == "" {
		return "", domain.ErrEmptyParameters
	}

	userId, err := u.tokenManager.Verify(token)
	if err != nil {
		u.logger.Error(err)
		return "", err
	}

	return userId, nil
}

func (u *Users) UserExists(userId string) (bool, error) {
	if userId == "" {
		return false, domain.ErrEmptyParameters
	}

	exists, err := u.repository.CheckExistence(userId)
	if err != nil {
		u.logger.Error(err)
		return false, err
	}

	return exists, nil
}
