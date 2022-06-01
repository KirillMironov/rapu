package service

import (
	"errors"
	"github.com/KirillMironov/rapu/users/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	usersRepository UsersRepository
	jwtService      JWTService
	logger          Logger
}

type UsersRepository interface {
	Create(domain.User) (userId string, err error)
	GetByEmail(email string) (domain.User, error)
	CheckExistence(userId string) (bool, error)
}

type JWTService interface {
	Generate(userId string) (token string, err error)
	Verify(token string) (userId string, err error)
}

type Logger interface {
	Error(args ...interface{})
}

func NewUsers(usersRepository UsersRepository, jwtService JWTService, logger Logger) *Users {
	return &Users{
		usersRepository: usersRepository,
		jwtService:      jwtService,
		logger:          logger,
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

	userId, err := u.usersRepository.Create(user)
	if err != nil {
		u.logger.Error(err)
		return "", err
	}

	token, err := u.jwtService.Generate(userId)
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

	user, err := u.usersRepository.GetByEmail(input.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", domain.ErrInvalidCredentials
		}
		u.logger.Error(err)
		return "", err
	}

	token, err := u.jwtService.Generate(user.Id)
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

	userId, err := u.jwtService.Verify(token)
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

	exists, err := u.usersRepository.CheckExistence(userId)
	if err != nil {
		u.logger.Error(err)
		return false, err
	}

	return exists, nil
}
