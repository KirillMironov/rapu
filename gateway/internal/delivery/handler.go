package delivery

import (
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/KirillMironov/rapu/gateway/pkg/echox"
	"github.com/KirillMironov/rapu/gateway/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	usersClient proto.UsersClient
	postsClient proto.PostsClient
	logger      logger.Logger
}

func NewHandler(usersClient proto.UsersClient, postsClient proto.PostsClient, logger logger.Logger) *Handler {
	return &Handler{
		usersClient: usersClient,
		postsClient: postsClient,
		logger:      logger,
	}
}

func (h Handler) InitRoutes() *echo.Echo {
	router := echo.New()
	router.Binder = echox.Binder{}
	router.Validator = echox.NewStructValidator()
	router.HTTPErrorHandler = echox.GRPCErrorHandler
	router.Use(
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderContentType, echo.HeaderContentLength, echo.HeaderAuthorization},
			AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		}),
	)

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/sign-up", h.signUp)
			users.POST("/sign-in", h.signIn)
		}
		posts := v1.Group("/posts")
		{
			posts.GET("/:userId", h.getPostsByUserId)
			posts.POST("", h.createPost, h.auth)
		}
	}

	return router
}
