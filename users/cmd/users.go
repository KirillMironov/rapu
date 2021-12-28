package main

import (
	"github.com/KirillMironov/rapu/config"
	"github.com/KirillMironov/rapu/internal/delivery"
	"github.com/KirillMironov/rapu/internal/repository/postgres"
	"github.com/KirillMironov/rapu/internal/service"
	"github.com/KirillMironov/rapu/pkg/auth"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net"
)

func main() {
	// Config
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Postgres
	db, err := sqlx.Connect("postgres", cfg.Postgres.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	// JWT manager
	tokenManager, err := auth.NewManager(cfg.Security.JWTKey, cfg.Security.TokenTTL)
	if err != nil {
		log.Fatal(err)
	}

	// App
	usersRepo := postgres.NewUsersRepository(db)
	usersService := service.NewUsersService(usersRepo, tokenManager)
	handler := delivery.NewHandler(usersService)

	listener, err := net.Listen("tcp", "localhost:"+cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server started on port %s", cfg.Port)
	log.Fatal(handler.Serve(listener))
}
