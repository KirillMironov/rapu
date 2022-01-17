package main

import (
	"github.com/KirillMironov/rapu/messenger/config"
	"github.com/KirillMironov/rapu/messenger/internal/delivery"
	repo "github.com/KirillMironov/rapu/messenger/internal/repository/redis"
	"github.com/KirillMironov/rapu/messenger/internal/service"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

func main() {
	// Config
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// App
	repository := repo.NewMessagesRepository(client)
	messagesService := service.NewMessagesService(repository)
	clientsService := service.NewClientsService(messagesService)
	handler := delivery.NewHandler(clientsService)

	log.Printf("messenger started on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler.InitRoutes()))
}
