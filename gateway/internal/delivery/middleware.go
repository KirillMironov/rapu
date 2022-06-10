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
	if c.Response().Committed {
		return
	}

	var httpError *echo.HTTPError

	switch v := err.(type) {
	case *echo.HTTPError:
		httpError = v
		if v.Internal != nil {
			if internalErr, ok := v.Internal.(*echo.HTTPError); ok {
				httpError = internalErr
			}
		}
	case interface{ GRPCStatus() *status.Status }:
		st, ok := status.FromError(err)
		if !ok {
			httpError = echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			break
		}

		switch st.Code() {
		case codes.InvalidArgument:
			httpError = echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		case codes.NotFound:
			httpError = echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		case codes.AlreadyExists:
			httpError = echo.NewHTTPError(http.StatusConflict, http.StatusText(http.StatusConflict))
		case codes.Unauthenticated:
			httpError = echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		default:
			h.logger.Error(err)
			httpError = echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	default:
		httpError = echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	code := httpError.Code
	message := httpError.Message
	if m, ok := httpError.Message.(string); ok {
		message = echo.Map{"message": m}
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(httpError.Code)
	} else {
		err = c.JSON(code, message)
	}
	if err != nil {
		h.logger.Error(err)
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
