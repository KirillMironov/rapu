package main

import (
	"context"
	"github.com/KirillMironov/rapu/posts/config"
	"github.com/KirillMironov/rapu/posts/internal/delivery"
	"github.com/KirillMironov/rapu/posts/internal/repository"
	"github.com/KirillMironov/rapu/posts/internal/service"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net"
	"os"
	"os/signal"
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

	// Mongo
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.ConnectionString))
	if err != nil {
		logger.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal(err)
	}

	db := client.Database(cfg.Mongo.DBName).Collection(cfg.Mongo.Collection)

	// App
	postsRepository := repository.NewPosts(db)
	postsService := service.NewPosts(postsRepository, cfg.MaxPostsPerPage)
	handler := delivery.NewHandler(postsService, logger)

	// gRPC Server
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("gRPC server started on port %s", cfg.Port)
	go func() {
		err = handler.Serve(listener)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Infof("shutting down gRPC server")
	handler.GracefulStop()
}
