package main

import (
	"grpc/internal/config"
	"grpc/internal/logger"
	"grpc/pkg/database"
	gRPCserver "grpc/pkg/gRPCserver"
	"grpc/pkg/shutdown"
	"log/slog"
)

func main() {
	cfg := config.LoadConfig()
	logger.Init(cfg.Env)
	slog.Info("Cfg launched successfully")
	db, err := database.DbInit(&cfg.Database)
	if err != nil {
		slog.Error("Database not initialization", "ERROR:", err.Error())
		return
	}
	go gRPCserver.Run(&cfg.Server, db)
	shutdown.Shutdown(db)
}
