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
	Error(args ...interface{})
}

func NewHandler(usersClient proto.UsersClient, postsClient proto.PostsClient, logger Logger) *Handler {
	return &Handler{
		usersClient: usersClient,
		postsClient: postsClient,
		logger:      logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), h.middleware)

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
			posts.Use(h.auth).POST("", h.createPost)
		}
	}

	return router
}
