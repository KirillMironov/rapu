package main

import (
	"context"
	"errors"
	"github.com/KirillMironov/rapu/gateway/config"
	"github.com/KirillMironov/rapu/gateway/internal/delivery"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	// Config
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	// gRPC Users client
	usersConn, err := grpc.Dial(cfg.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	defer usersConn.Close()

	usersClient := proto.NewUsersClient(usersConn)

	// gRPC Posts client
	postsConn, err := grpc.Dial(cfg.PostsServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	defer postsConn.Close()

	postsClient := proto.NewPostsClient(postsConn)

	// App
	handler := delivery.NewHandler(usersClient, postsClient, logger)

	// Gin
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.InitRoutes(),
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	logger.Infof("gateway started on port %s", cfg.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	logger.Info("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
}
