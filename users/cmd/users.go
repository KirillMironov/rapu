package main

import (
	"github.com/KirillMironov/rapu/users/config"
	"github.com/KirillMironov/rapu/users/internal/delivery"
	"github.com/KirillMironov/rapu/users/internal/repository/postgres"
	"github.com/KirillMironov/rapu/users/internal/service"
	"github.com/KirillMironov/rapu/users/pkg/auth"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net"
	"time"
)

func main() {
	// Logger
	zapCfg := zap.NewProductionConfig()
	zapCfg.Encoding = "console"
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapCfg.EncoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("[" + time.Format("Jan 2 15:04:05.000") + "]")
	}
	zapCfg.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("(" + caller.TrimmedPath() + ")")
	}

	zapLogger, err := zapCfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	defer logger.Sync()

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
	tokenManager, err := auth.NewManager(cfg.Security.JWTKey, cfg.Security.TokenTTL)
	if err != nil {
		logger.Fatal(err)
	}

	// App
	usersRepo := postgres.NewUsersRepository(db)
	usersService := service.NewUsersService(usersRepo, tokenManager)
	handler := delivery.NewHandler(usersService, logger)

	listener, err := net.Listen("tcp", "localhost:"+cfg.Port)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("server started on port %s", cfg.Port)
	logger.Fatal(handler.Serve(listener))
}
