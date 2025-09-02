package main

import (
	"grpc/internal/config"
	"grpc/internal/logger"
	"grpc/pkg/server"
	"grpc/pkg/shutdown"
	"log/slog"
)

func main() {
	cfg := config.LoadConfig()
	logger.Init(cfg.Env)
	slog.Info("Cfg and logger launched successfully")

	go server.Run(&cfg.Server)

	shutdown.Shutdown()
}
