package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const accessTokenKey = "access_token"

type signUpForm struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var form signUpForm

	err := c.BindJSON(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	resp, err := h.usersClient.SignUp(context.Background(), &proto.SignUpRequest{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
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
		case codes.AlreadyExists:
			c.Status(http.StatusUnauthorized)
			return
		default:
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{accessTokenKey: resp.GetAccessToken()})
}

type signInForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var form signInForm

	err := c.BindJSON(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	resp, err := h.usersClient.SignIn(context.Background(), &proto.SignInRequest{
		Email:    form.Email,
		Password: form.Password,
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
		case codes.Unauthenticated:
			c.Status(http.StatusUnauthorized)
			return
		default:
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{accessTokenKey: resp.GetAccessToken()})
}
