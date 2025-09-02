package main

import (
	"grpc/internal/config"
	"grpc/pkg/server"
	"grpc/pkg/shutdown"
)

func main() {
	cfg := config.LoadConfig().Server

	go server.Run(&cfg)

	shutdown.Shutdown()
}
