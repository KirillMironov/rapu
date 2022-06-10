package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/labstack/echo/v4"
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
