package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createPostForm struct {
	UserId  string `json:"user_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (h *Handler) createPost(c *gin.Context) {
	var form createPostForm

	err := c.BindJSON(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	_, err = h.postsClient.Create(context.Background(), &proto.CreateRequest{
		UserId:  form.UserId,
		Message: form.Message,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.Status(http.StatusCreated)
}

type getPostsByUserIdForm struct {
	Offset string `form:"offset"`
	Limit  int64  `form:"limit"`
}

func (h *Handler) getPostsByUserId(c *gin.Context) {
	var form getPostsByUserIdForm

	err := c.Bind(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	resp, err := h.postsClient.GetByUserId(context.Background(), &proto.GetByUserIdRequest{
		UserId: c.Param("userId"),
		Offset: form.Offset,
		Limit:  form.Limit,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusOK, resp.GetPosts())
}
