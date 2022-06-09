package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/messenger/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()
	router.Validator = &Validator{validator: validator.New()}
	router.Binder = &Binder{}
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
			messenger.GET("/connect", h.connect, h.auth)
		}
	}

	return router
}

func (h *Handler) connect(c echo.Context) error {
	var form struct {
		ToUserId string `form:"toUserId" validate:"required"`
	}

	err := c.Bind(&form)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var client = domain.Client{
		ToUserId: form.ToUserId,
		Conn:     conn,
	}

	err = h.clientsService.Connect(c.Request().Context(), c.Get(jwtKey).(string), client)
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

	return nil
}
