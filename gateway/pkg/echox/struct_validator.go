package echox

import "github.com/go-playground/validator/v10"

type structValidator struct {
	validator *validator.Validate
}

func NewStructValidator() *structValidator {
	return &structValidator{validator.New()}
}

func (sv structValidator) Validate(i interface{}) error {
	return sv.validator.Struct(i)
}
