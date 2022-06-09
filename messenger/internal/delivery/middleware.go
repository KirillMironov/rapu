package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const jwtKey = "jwt"

func (h Handler) auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(echo.HeaderAuthorization)
		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		token := strings.Split(header, "Bearer ")
		if len(token) != 2 || token[1] == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		c.Set(jwtKey, token[1])
		return next(c)
	}
}

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
