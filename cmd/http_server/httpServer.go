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

	serv := grpcconect.NewServer(cfg.Http_server.Grpc_client.Host, cfg.Http_server.Grpc_client.Port)

	go httpServer.Run(serv, cfg.Http_server.Host, cfg.Http_server.Port)
	httpshutdown.Shutdown(serv)
}
