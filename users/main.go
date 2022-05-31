package main

import (
	"github.com/KirillMironov/rapu/users/config"
	"github.com/KirillMironov/rapu/users/internal/delivery"
	"github.com/KirillMironov/rapu/users/internal/repository"
	"github.com/KirillMironov/rapu/users/internal/service"
	"github.com/KirillMironov/rapu/users/pkg/jwt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net"
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

	// JWT manager
	tokenManager, err := jwt.NewTokenManager(cfg.Security.JWTKey, cfg.Security.TokenTTL)
	if err != nil {
		logger.Fatal(err)
	}

	// App
	usersRepository := repository.NewUsers(db)
	usersService := service.NewUsers(usersRepository, tokenManager, logger)
	handler := delivery.NewHandler(usersService)

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("started on port %s", cfg.Port)
	logger.Fatal(handler.Serve(listener))
}
