package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/KirillMironov/rapu/gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	accessToken = "access_token"
	userId      = "user_id"
)

type Handler struct {
	client proto.UsersClient
	logger logger.Logger
}

func NewHandler(client proto.UsersClient, logger logger.Logger) *Handler {
	return &Handler{client: client, logger: logger}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/sign-up", h.signUp)
			users.POST("/sign-in", h.signIn)
			users.GET("/auth", h.auth)
		}
	}

	return r
}

func (h *Handler) signUp(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.BindJSON(&credentials)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	resp, err := h.client.SignUp(context.Background(), &proto.SignUpRequest{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]string{accessToken: resp.AccessToken})
}

func (h *Handler) signIn(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.BindJSON(&credentials)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	resp, err := h.client.SignIn(context.Background(), &proto.SignInRequest{
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{accessToken: resp.AccessToken})
}

func (h *Handler) auth(c *gin.Context) {
	token, ok := c.GetQuery("token")
	if !ok {
		h.logger.Info(http.StatusBadRequest)
		return
	}

	resp, err := h.client.Authenticate(context.Background(), &proto.AuthRequest{
		AccessToken: token,
	})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{userId: resp.UserId})
}
