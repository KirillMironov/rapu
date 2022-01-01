package main

import (
	"github.com/KirillMironov/rapu/gateway/config"
	"github.com/KirillMironov/rapu/gateway/internal/delivery"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	// Config
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	// GRPC Client
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(cfg.ServerAddress, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewUsersClient(conn)

	// App
	handler := delivery.NewHandler(client)
	log.Printf("gateway started on port %s", cfg.Port)
	log.Fatal(handler.InitRoutes().Run(":" + cfg.Port))
}
