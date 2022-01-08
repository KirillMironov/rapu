package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	accessToken = "access_token"
	userId      = "user_id"
)

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
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{accessToken: resp.GetAccessToken()})
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
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{accessToken: resp.GetAccessToken()})
}

type authForm struct {
	AccessToken string `json:"access_token" binding:"required"`
}

func (h *Handler) auth(c *gin.Context) {
	var form authForm

	err := c.BindJSON(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info("access token was not provided")
		return
	}

	resp, err := h.usersClient.Authenticate(context.Background(), &proto.AuthRequest{
		AccessToken: form.AccessToken,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{userId: resp.GetUserId()})
}
