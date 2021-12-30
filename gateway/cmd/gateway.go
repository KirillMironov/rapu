package main

import (
	"github.com/KirillMironov/rapu/gateway/internal/delivery"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const (
	port          = "7002"
	serverAddress = "localhost:7001"
)

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewUsersClient(conn)

	handler := delivery.NewHandler(client)
	log.Printf("gateway started on port %s", port)
	log.Fatal(handler.InitRoutes().Run(":" + port))
}
