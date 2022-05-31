package main

import (
	"context"
	"github.com/KirillMironov/rapu/posts/config"
	"github.com/KirillMironov/rapu/posts/internal/delivery"
	_repository "github.com/KirillMironov/rapu/posts/internal/repository/mongo"
	_service "github.com/KirillMironov/rapu/posts/internal/service"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net"
)

var ctx = context.Background()

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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.ConnectionString))
	if err != nil {
		logger.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal()
	}

	db := client.Database(cfg.Mongo.DBName).Collection(cfg.Mongo.Collection)

	// App
	repository := _repository.NewPostsRepository(db)
	service := _service.NewPostsService(repository, cfg.MaxPostsPerPage)
	handler := delivery.NewHandler(service, logger)

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("server started on port %s", cfg.Port)
	logger.Fatal(handler.Serve(listener))
}