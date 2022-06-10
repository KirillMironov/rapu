package echox

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	Validator *validator.Validate
}

func (sv StructValidator) Validate(i interface{}) error {
	return sv.Validator.Struct(i)
}
