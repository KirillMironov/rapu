package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createPostRequest struct {
	UserId  string `json:"user_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type getByUserIdRequest struct {
	UserId string `form:"userId" binding:"required"`
	Offset string `form:"offset" binding:"required"`
	Limit  int64  `form:"limit" binding:"required"`
}

func (h *Handler) create(c *gin.Context) {
	var req createPostRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	_, err = h.postsClient.Create(context.Background(), &proto.CreateRequest{
		UserId:  req.UserId,
		Message: req.Message,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}
}

func (h *Handler) getByUserId(c *gin.Context) {
	var req getByUserIdRequest

	err := c.Bind(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	resp, err := h.postsClient.GetByUserId(context.Background(), &proto.GetByUserIdRequest{
		UserId: req.UserId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusOK, resp.Posts)
}
