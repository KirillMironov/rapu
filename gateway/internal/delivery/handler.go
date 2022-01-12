package delivery

import (
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	usersClient proto.UsersClient
	postsClient proto.PostsClient
	logger      Logger
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

func NewHandler(client proto.UsersClient, postsClient proto.PostsClient, logger Logger) *Handler {
	return &Handler{
		usersClient: client,
		postsClient: postsClient,
		logger:      logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/sign-up", h.signUp)
			users.POST("/sign-in", h.signIn)
			users.POST("/auth", h.auth)
		}

		posts := v1.Group("/posts")
		{
			posts.POST("", h.createPost)
			posts.GET("/:userId", h.getPostsByUserId)
		}
	}

	return router
}
