package delivery

import (
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/labstack/echo/v4"
	"strings"
)

const userIdKey = "userId"

func (h Handler) auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			return echo.ErrUnauthorized
		}

		token := strings.Split(header, "Bearer ")
		if len(token) != 2 || token[1] == "" {
			return echo.ErrUnauthorized
		}

		resp, err := h.usersClient.Authenticate(c.Request().Context(), &proto.AuthRequest{AccessToken: token[1]})
		if err != nil {
			return echo.ErrUnauthorized
		}

		userId := resp.GetUserId()
		if userId == "" {
			return echo.ErrUnauthorized
		}

		c.Set(userIdKey, userId)
		return next(c)
	}
}
