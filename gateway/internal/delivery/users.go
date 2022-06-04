package delivery

import (
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const accessTokenKey = "access_token"

func (h *Handler) signUp(c *gin.Context) {
	var form struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.BindJSON(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := h.usersClient.SignUp(c, &proto.SignUpRequest{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
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
		case codes.AlreadyExists:
			c.Status(http.StatusConflict)
		default:
			h.logger.Error(err)
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{accessTokenKey: resp.GetAccessToken()})
}

func (h *Handler) signIn(c *gin.Context) {
	var form struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.BindJSON(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := h.usersClient.SignIn(c, &proto.SignInRequest{
		Email:    form.Email,
		Password: form.Password,
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
		case codes.Unauthenticated:
			c.Status(http.StatusUnauthorized)
		default:
			h.logger.Error(err)
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{accessTokenKey: resp.GetAccessToken()})
}
