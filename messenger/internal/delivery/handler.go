package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/messenger/domain"
	"github.com/KirillMironov/rapu/messenger/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	AccessToken string `form:"access_token" binding:"required"`
}

func (h *Handler) connect(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	var form connectForm
	err = c.Bind(&form)
	if err != nil {
		cm := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "not enough query parameters")
		_ = conn.WriteMessage(websocket.CloseMessage, cm)
		conn.Close()
		return
	}

	authResp, err := h.usersClient.Authenticate(context.Background(), &proto.AuthRequest{AccessToken: form.AccessToken})
	if err != nil || authResp.GetUserId() == "" {
		defer conn.Close()
		st, ok := status.FromError(err)
		if !ok {
			_ = conn.WriteMessage(websocket.CloseInternalServerErr, nil)
			h.logger.Error(err)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			cm := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, st.Message())
			_ = conn.WriteMessage(websocket.CloseMessage, cm)
		case codes.Unauthenticated:
			cm := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, st.Message())
			_ = conn.WriteMessage(websocket.CloseMessage, cm)
		default:
			_ = conn.WriteMessage(websocket.CloseInternalServerErr, nil)
			h.logger.Error(err)
		}
		return
	}

	resp, err := h.usersClient.UserExists(context.Background(), &proto.UserExistsRequest{UserId: form.ToUserId})
	if err != nil || resp.GetExists() == false {
		defer conn.Close()
		st, ok := status.FromError(err)
		if !ok {
			_ = conn.WriteMessage(websocket.CloseInternalServerErr, nil)
			h.logger.Error(err)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			cm := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, st.Message())
			_ = conn.WriteMessage(websocket.CloseMessage, cm)
		case codes.NotFound:
			cm := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, st.Message())
			_ = conn.WriteMessage(websocket.CloseMessage, cm)
		default:
			_ = conn.WriteMessage(websocket.CloseInternalServerErr, nil)
			h.logger.Error(err)
		}
		return
	}

	var client = domain.Client{
		UserId:   authResp.GetUserId(),
		ToUserId: form.ToUserId,
		Conn:     conn,
	}

	h.service.Connect(client)
}
