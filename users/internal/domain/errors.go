package domain

import "errors"

var (
	ErrEmptyParameters    = errors.New("received one or more empty parameters")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
