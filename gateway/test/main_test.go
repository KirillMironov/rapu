package test

import (
	"github.com/KirillMironov/rapu/gateway/internal/delivery"
	"github.com/KirillMironov/rapu/gateway/test/mocks"
	"github.com/gin-gonic/gin"
	"testing"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	handler := delivery.NewHandler(mocks.UsersClientMock{}, mocks.LoggerMock{})
	router = handler.InitRoutes()

	m.Run()
}
