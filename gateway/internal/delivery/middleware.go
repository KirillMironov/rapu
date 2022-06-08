package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

const userIdKey = "userId"

func (h Handler) auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		token := strings.Split(header, "Bearer ")
		if len(token) != 2 || token[1] == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		resp, err := h.usersClient.Authenticate(context.Background(), &proto.AuthRequest{AccessToken: token[1]})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		userId := resp.GetUserId()
		if userId == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		c.Set(userIdKey, userId)
		return next(c)
	}
}

func (h Handler) errorHandler(err error, c echo.Context) {
	var (
		code    int
		message interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	} else if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			code = http.StatusBadRequest
		case codes.NotFound:
			code = http.StatusNotFound
		case codes.AlreadyExists:
			code = http.StatusConflict
		default:
			h.logger.Error(err)
			code = http.StatusInternalServerError
		}
	} else {
		h.logger.Error(err)
		code = http.StatusInternalServerError
		message = http.StatusText(code)
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				h.logger.Error(err)
			}
		} else {
			err := c.JSON(code, echo.Map{"message": message})
			if err != nil {
				h.logger.Error(err)
			}
		}
	}
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
