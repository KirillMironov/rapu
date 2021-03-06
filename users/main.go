package main

import (
	"github.com/KirillMironov/rapu/users/config"
	"github.com/KirillMironov/rapu/users/internal/delivery"
	"github.com/KirillMironov/rapu/users/internal/repository"
	"github.com/KirillMironov/rapu/users/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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

	// Postgres
	db, err := sqlx.Connect("postgres", cfg.Postgres.ConnectionString)
	if err != nil {
		logger.Fatal(err)
	}

	// App
	usersRepository := repository.NewUsers(db)
	jwtService, err := service.NewJWT(cfg.Security.JWTKey, cfg.Security.TokenTTL)
	if err != nil {
		logger.Fatal(err)
	}
	usersService := service.NewUsers(usersRepository, jwtService)
	handler := delivery.NewHandler(usersService, logger)

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
