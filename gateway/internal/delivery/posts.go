package delivery

import (
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (h *Handler) createPost(c *gin.Context) {
	var form struct {
		Message string `json:"message" binding:"required"`
	}

	err := c.BindJSON(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, err = h.postsClient.Create(c, &proto.CreateRequest{
		UserId:  c.GetString(userIdKey),
		Message: form.Message,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			h.logger.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			c.Status(http.StatusBadRequest)
		default:
			h.logger.Error(err)
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) getPostsByUserId(c *gin.Context) {
	var form struct {
		Offset string `form:"offset"`
		Limit  int64  `form:"limit"`
	}

	err := c.Bind(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := h.postsClient.GetByUserId(c, &proto.GetByUserIdRequest{
		UserId: c.Param("userId"),
		Offset: form.Offset,
		Limit:  form.Limit,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			h.logger.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			c.Status(http.StatusBadRequest)
		case codes.NotFound:
			c.Status(http.StatusNotFound)
		default:
			h.logger.Error(err)
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, resp.GetPosts())
}
