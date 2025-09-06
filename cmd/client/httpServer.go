package main

import (
	"grpc/internal/config"
	"grpc/internal/logger"
	"grpc/pkg/httpServer"
	grpcconect "grpc/pkg/httpServer/grpcConect"
	httpshutdown "grpc/pkg/shutdown/httpServer"
	"log/slog"
)

func main() {
	cfg := config.LoadConfig()
	logger.Init(cfg.Env)
	slog.Info("Cfg launched successfully")

	serv := grpcconect.NewServer(cfg.Grpc_server.Host, cfg.Grpc_server.Port)

	go httpServer.Run(serv, cfg.Http_server.Host, cfg.Http_server.Port)
	httpshutdown.Shutdown(serv)
}
