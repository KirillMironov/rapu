package main

import (
	"github.com/KirillMironov/rapu/messenger/config"
	"github.com/KirillMironov/rapu/messenger/internal/delivery"
	"github.com/KirillMironov/rapu/messenger/internal/delivery/proto"
	repo "github.com/KirillMironov/rapu/messenger/internal/repository/redis"
	"github.com/KirillMironov/rapu/messenger/internal/service"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
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
	bus := repo.NewMessagesBus(client)
	repository := repo.NewMessagesRepository(client)
	messagesService := service.NewMessagesService(bus, repository, logger)
	clientsService := service.NewClientsService(messagesService)
	handler := delivery.NewHandler(usersClient, clientsService, logger)

	logger.Infof("messenger started on port %s", cfg.Port)
	logger.Fatal(http.ListenAndServe(":"+cfg.Port, handler.InitRoutes()))
}
