package main

import (
	"grpc/internal/config"
	"grpc/internal/logger"
	"grpc/pkg/httpServer"
	"log/slog"
)

func main() {
	cfg := config.LoadConfig()
	logger.Init(cfg.Env)
	slog.Info("Cfg launched successfully")
	httpServer.Run()
}
