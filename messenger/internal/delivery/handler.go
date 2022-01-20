package delivery

import (
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
		messenger := v1.Group("messenger").Use(h.auth)
		{
			messenger.GET("/connect", h.connect)
		}
	}

	return router
}

func (h *Handler) connect(c *gin.Context) {
	toUserId, ok := c.GetQuery("toUserId")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	var client = domain.Client{
		UserId:   c.GetString(userIdKey),
		ToUserId: toUserId,
		Conn:     conn,
	}

	h.service.Connect(client)
}
