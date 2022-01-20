package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/messenger/internal/delivery/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

const userIdKey = "userId"

func (h *Handler) middleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")
	c.Header("Content-Type", "application/json")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
}

func (h *Handler) auth(c *gin.Context) {
	accessToken, ok := c.GetQuery("access_token")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		h.logger.Info("access token was not provided")
		return
	}

	resp, err := h.usersClient.Authenticate(context.Background(), &proto.AuthRequest{AccessToken: accessToken})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		h.logger.Info(err)
		return
	}

	userId := resp.GetUserId()
	if userId == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		h.logger.Error("missing userId")
		return
	}

	c.Set(userIdKey, userId)
	c.Next()
}
