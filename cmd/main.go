package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"grpc-auth/internal/config"
	customLogger "grpc-auth/internal/logger"
	pb "grpc-auth/internal/protos"
	"grpc-auth/internal/repo"
	"grpc-auth/internal/service"
)

func main() {
	// Загружаем конфигурацию из переменных окружения
	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	// Инициализация логгера
	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	// Подключение к PostgreSQL
	repository, err := repo.NewRepository(context.Background(), cfg.PostgreSQL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize repository"))
	}

	// Инициализация API
	grpcServer := grpc.NewServer()
	authService := service.NewAuthService(repository, logger)
	pb.RegisterAuthServiceServer(grpcServer, authService)

	listen, err := net.Listen("tcp", cfg.Grpc.Port)
	if err != nil {
		logger.Fatal(err, "error creating listener")
	}

	go func() {
		logger.Infof("grpc started listening on port: %s", cfg.Grpc.Port)
		if err := grpcServer.Serve(listen); err != nil {
			logger.Fatal("failed to serve", err)
		}
	}()

	quitSig := make(chan os.Signal, 1)
	signal.Notify(quitSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-quitSig
	logger.Infof("Shutting down grpc server gracefully...")
	grpcServer.GracefulStop()

	if err := repository.ShuttingDownPostgres(); err != nil {
		log.Fatal("Error closing connection", zap.Error(err))
	}

	fmt.Println("server stopped")

}
