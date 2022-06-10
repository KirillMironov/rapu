package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	validator *validator.Validate
}

func (v Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

type Binder struct{}

func (Binder) Bind(i interface{}, c echo.Context) error {
	var binder echo.DefaultBinder

	err := binder.Bind(i, c)
	if err != nil {
		return err
	}

	return c.Validate(i)
}
