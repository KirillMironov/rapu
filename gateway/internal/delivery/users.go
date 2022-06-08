package delivery

import (
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/labstack/echo/v4"
	"net/http"
)

const accessTokenKey = "access_token"

func (h Handler) signUp(c echo.Context) error {
	var form struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.usersClient.SignUp(c.Request().Context(), &proto.SignUpRequest{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{accessTokenKey: resp.GetAccessToken()})
}

func (h Handler) signIn(c echo.Context) error {
	var form struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.usersClient.SignIn(c.Request().Context(), &proto.SignInRequest{
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{accessTokenKey: resp.GetAccessToken()})
}
