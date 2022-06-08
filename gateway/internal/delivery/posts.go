package delivery

import (
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) createPost(c echo.Context) error {
	var form struct {
		Message string `json:"message" validate:"required"`
	}

	if err := c.Bind(&form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := h.postsClient.Create(c.Request().Context(), &proto.CreateRequest{
		UserId:  c.Get(userIdKey).(string),
		Message: form.Message,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h Handler) getPostsByUserId(c echo.Context) error {
	var form struct {
		Offset string `form:"offset"`
		Limit  int64  `form:"limit"`
	}

	if err := c.Bind(&form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.postsClient.GetByUserId(c.Request().Context(), &proto.GetByUserIdRequest{
		UserId: c.Param("userId"),
		Offset: form.Offset,
		Limit:  form.Limit,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp.GetPosts())
}
