package service

import (
	"context"
	"errors"
	"github.com/KirillMironov/rapu/users/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	usersRepository UsersRepository
	jwtService      JWTService
}

type UsersRepository interface {
	Create(context.Context, domain.User) (userId string, err error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	CheckExistence(ctx context.Context, userId string) (bool, error)
}

type JWTService interface {
	Generate(userId string) (token string, err error)
	Verify(token string) (userId string, err error)
}

func NewUsers(usersRepository UsersRepository, jwtService JWTService) *Users {
	return &Users{
		usersRepository: usersRepository,
		jwtService:      jwtService,
	}
}

func (u Users) SignUp(ctx context.Context, user domain.User) (string, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return "", domain.ErrEmptyParameters
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hash)

	userId, err := u.usersRepository.Create(ctx, user)
	if err != nil {
		return "", err
	}

	token, err := u.jwtService.Generate(userId)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u Users) SignIn(ctx context.Context, input domain.User) (string, error) {
	if input.Email == "" || input.Password == "" {
		return "", domain.ErrEmptyParameters
	}

	user, err := u.usersRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	token, err := u.jwtService.Generate(user.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u Users) Authenticate(token string) (string, error) {
	if token == "" {
		return "", domain.ErrEmptyParameters
	}

	userId, err := u.jwtService.Verify(token)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (u Users) UserExists(ctx context.Context, userId string) (bool, error) {
	if userId == "" {
		return false, domain.ErrEmptyParameters
	}

	exists, err := u.usersRepository.CheckExistence(ctx, userId)
	if err != nil {
		return false, err
	}

	return exists, nil
}
