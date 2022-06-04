package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/messenger/internal/domain"
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
	clientsService ClientsService
	logger         Logger
}

type ClientsService interface {
	Connect(ctx context.Context, accessToken string, client domain.Client) error
}

type Logger interface {
	Error(args ...interface{})
}

func NewHandler(clientsService ClientsService, logger Logger) *Handler {
	return &Handler{
		clientsService: clientsService,
		logger:         logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), h.middleware)

	v1 := router.Group("/api/v1")
	{
		messenger := v1.Group("/messenger")
		{
			messenger.GET("/connect", h.connect)
		}
	}

	return router
}

func (h *Handler) connect(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	var form struct {
		ToUserId    string `form:"toUserId" binding:"required"`
		AccessToken string `form:"access_token" binding:"required"`
	}

	err = c.Bind(&form)
	if err != nil {
		closeMessage := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error())
		_ = conn.WriteMessage(websocket.CloseMessage, closeMessage)
		conn.Close()
		return
	}

	var client = domain.Client{
		ToUserId: form.ToUserId,
		Conn:     conn,
	}

	err = h.clientsService.Connect(c, form.AccessToken, client)
	if err != nil {
		var closeMessage []byte
		switch err {
		case domain.ErrEmptyParameters:
			closeMessage = websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error())
		case domain.ErrUserNotFound:
			closeMessage = websocket.FormatCloseMessage(websocket.ClosePolicyViolation, err.Error())
		default:
			closeMessage = websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())
		}
		_ = conn.WriteMessage(websocket.CloseMessage, closeMessage)
		conn.Close()
	}
}
