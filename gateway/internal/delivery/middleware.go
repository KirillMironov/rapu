package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const userIdKey = "userId"

func (h *Handler) middleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")
	c.Header("Content-Type", "application/json")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
}

func (h *Handler) auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := strings.Split(header, "Bearer ")
	if len(token) != 2 || token[1] == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	resp, err := h.usersClient.Authenticate(context.Background(), &proto.AuthRequest{AccessToken: token[1]})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId := resp.GetUserId()
	if userId == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(userIdKey, userId)
	c.Next()
}
