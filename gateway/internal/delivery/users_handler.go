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

type signUpCredentials struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signInCredentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var credentials signUpCredentials

	err := c.BindJSON(&credentials)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	resp, err := h.usersClient.SignUp(context.Background(), &proto.SignUpRequest{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{accessToken: resp.GetAccessToken()})
}

func (h *Handler) signIn(c *gin.Context) {
	var credentials signInCredentials

	err := c.BindJSON(&credentials)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	resp, err := h.usersClient.SignIn(context.Background(), &proto.SignInRequest{
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{accessToken: resp.GetAccessToken()})
}

func (h *Handler) auth(c *gin.Context) {
	var req authRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info("access token was not provided")
		return
	}

	resp, err := h.usersClient.Authenticate(context.Background(), &proto.AuthRequest{
		AccessToken: req.AccessToken,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{userId: resp.GetUserId()})
}
