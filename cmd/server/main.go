package main

import (
	"grpc/internal/config"
	"grpc/internal/server"
)

func main() {
	cfg := config.LoadConfig().Server

	server.Run(&cfg)

}
