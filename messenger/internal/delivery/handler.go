package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/messenger/domain"
	"github.com/KirillMironov/rapu/messenger/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Handler struct {
	usersClient proto.UsersClient
	service     domain.ClientsService
	logger      Logger
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

func NewHandler(usersClient proto.UsersClient, service domain.ClientsService, logger Logger) *Handler {
	return &Handler{
		usersClient: usersClient,
		service:     service,
		logger:      logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), h.middleware)

	v1 := router.Group("/api/v1")
	{
		messenger := v1.Group("messenger")
		{
			messenger.GET("/connect", h.connect)
		}
	}

	return router
}

type connectForm struct {
	ToUserId    string `form:"toUserId" binding:"required"`
	AccessToken string `form:"accessToken" binding:"required"`
}

func (h *Handler) connect(c *gin.Context) {
	var form connectForm

	err := c.Bind(&form)
	if err != nil {
		c.Status(http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	resp, err := h.usersClient.Authenticate(context.Background(), &proto.AuthRequest{AccessToken: form.AccessToken})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	userId := resp.GetUserId()
	if userId == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	var client = domain.Client{
		UserId:   userId,
		ToUserId: form.ToUserId,
		Conn:     conn,
	}

	h.service.Connect(client)
}
