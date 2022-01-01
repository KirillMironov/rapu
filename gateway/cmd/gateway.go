package main

import (
	"github.com/KirillMironov/rapu/gateway/config"
	"github.com/KirillMironov/rapu/gateway/internal/delivery"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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

	// GRPC Client
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(cfg.ServerAddress, opts...)
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewUsersClient(conn)

	// App
	handler := delivery.NewHandler(client, logger)
	logger.Infof("gateway started on port %s", cfg.Port)
	logger.Fatal(handler.InitRoutes().Run(":" + cfg.Port))
}
