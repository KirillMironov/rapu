package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/pkg/echox"
	"github.com/KirillMironov/rapu/gateway/pkg/logger"
	"github.com/KirillMironov/rapu/messenger/internal/domain"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Handler struct {
	clientsService ClientsService
	logger         logger.Logger
}

type ClientsService interface {
	Connect(ctx context.Context, accessToken string, client domain.Client) error
}

func NewHandler(clientsService ClientsService, logger logger.Logger) *Handler {
	return &Handler{
		clientsService: clientsService,
		logger:         logger,
	}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()
	router.Binder = echox.Binder{}
	router.Validator = echox.NewStructValidator()
	router.Use(
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderContentType, echo.HeaderContentLength, echo.HeaderAuthorization},
			AllowMethods: []string{echo.GET, echo.OPTIONS},
		}),
	)

	v1 := router.Group("/api/v1")
	{
		messenger := v1.Group("/messenger")
		{
			messenger.GET("/connect", h.connect)
		}
	}

	return router
}

func (h *Handler) connect(c echo.Context) error {
	var form struct {
		ToUserId string `query:"toUserId" validate:"required"`
	}

	err := c.Bind(&form)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	messageType, accessToken, err := conn.ReadMessage()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if messageType != websocket.TextMessage {
		return echo.NewHTTPError(http.StatusBadRequest, "unsupported message type")
	}

	var client = domain.Client{
		ToUserId: form.ToUserId,
		Conn:     conn,
	}

	err = h.clientsService.Connect(c.Request().Context(), string(accessToken), client)
	if err != nil {
		var closeMessage []byte
		switch err {
		case domain.ErrEmptyParameters:
			closeMessage = websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error())
		case domain.ErrUserNotFound:
			closeMessage = websocket.FormatCloseMessage(websocket.ClosePolicyViolation, err.Error())
		default:
			h.logger.Error(err)
			closeMessage = websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())
		}
		_ = conn.WriteMessage(websocket.CloseMessage, closeMessage)
		conn.Close()
	}

	return nil
}
