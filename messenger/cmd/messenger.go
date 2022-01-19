package main

import (
	"github.com/KirillMironov/rapu/messenger/config"
	"github.com/KirillMironov/rapu/messenger/internal/delivery"
	repo "github.com/KirillMironov/rapu/messenger/internal/repository/redis"
	"github.com/KirillMironov/rapu/messenger/internal/service"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
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

	// App
	bus := repo.NewMessagesBus(client)
	repository := repo.NewMessagesRepository(client)
	messagesService := service.NewMessagesService(bus, repository, logger)
	clientsService := service.NewClientsService(messagesService)
	handler := delivery.NewHandler(clientsService, logger)

	logger.Infof("messenger started on port %s", cfg.Port)
	logger.Fatal(http.ListenAndServe(":"+cfg.Port, handler.InitRoutes()))
}
