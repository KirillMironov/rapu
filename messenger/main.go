package main

import (
	"context"
	"errors"
	"github.com/KirillMironov/rapu/messenger/config"
	"github.com/KirillMironov/rapu/messenger/internal/delivery"
	"github.com/KirillMironov/rapu/messenger/internal/delivery/proto"
	"github.com/KirillMironov/rapu/messenger/internal/repository"
	"github.com/KirillMironov/rapu/messenger/internal/service"
	"github.com/go-redis/redis"
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

	// Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer client.Close()

	err = client.Ping().Err()
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

	// App
	messagesBus := repository.NewMessagesBus(client)
	messagesRepository := repository.NewMessages(client)
	messagesService := service.NewMessages(messagesBus, messagesRepository, logger)
	clientsService := service.NewClients(usersClient, messagesService, logger)
	handler := delivery.NewHandler(clientsService, logger)

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

	logger.Infof("messenger started on port %s", cfg.Port)

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
