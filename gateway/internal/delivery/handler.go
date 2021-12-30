package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type (
	signUpCredentials struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	signInCredentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)

type Handler struct {
	client proto.UsersClient
}

func NewHandler(client proto.UsersClient) *Handler {
	return &Handler{client: client}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

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
	var credentials signUpCredentials
	err := c.BindJSON(&credentials)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Println(err)
		return
	}

	resp, err := h.client.SignUp(context.Background(), &proto.SignUpRequest{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusUnauthorized)
	}

	c.JSON(http.StatusCreated, resp.AccessToken)
}

func (h *Handler) signIn(c *gin.Context) {
	var credentials signInCredentials
	err := c.BindJSON(&credentials)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Println(err)
		return
	}

	resp, err := h.client.SignIn(context.Background(), &proto.SignInRequest{
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusUnauthorized)
	}

	c.JSON(http.StatusOK, resp.AccessToken)
}

func (h *Handler) auth(c *gin.Context) {
	token, ok := c.GetQuery("token")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := h.client.Authenticate(context.Background(), &proto.AuthRequest{
		AccessToken: token,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusUnauthorized)
	}

	c.JSON(http.StatusOK, resp.UserId)
}
