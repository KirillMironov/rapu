package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		st, ok := status.FromError(err)
		if !ok {
			c.Status(http.StatusInternalServerError)
			h.logger.Error(err)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			c.Status(http.StatusBadRequest)
			return
		default:
			c.Status(http.StatusInternalServerError)
			return
		}
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
		st, ok := status.FromError(err)
		if !ok {
			c.Status(http.StatusInternalServerError)
			h.logger.Error(err)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			c.Status(http.StatusBadRequest)
			return
		case codes.NotFound:
			c.Status(http.StatusNotFound)
			return
		default:
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.String(http.StatusOK, string(resp.GetPosts()))
}
