package domain

import "errors"

var (
	ErrEmptyParameters = errors.New("received one or more empty parameters")
	ErrEmptyResult     = errors.New("empty result")
)
