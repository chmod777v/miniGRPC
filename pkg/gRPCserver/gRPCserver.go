package gRPCserver

import (
	"grpc/internal/config"
	my_grpc "grpc/internal/grpc"
	"grpc/pkg/database"
	g_serv "grpc/pkg/proto"
	"log/slog"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(s *config.Grpc_server, db *database.Database) {
	addr := s.Host + ":" + strconv.Itoa(s.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic("Failed to listen:" + err.Error())
	}

	ser := grpc.NewServer()
	reflection.Register(ser)

	g_serv.RegisterServiceServer(ser, &my_grpc.Server{Db: db})
	slog.Info("Server listening", "Host", lis.Addr())

	if err = ser.Serve(lis); err != nil {
		panic("Failed to serve:" + err.Error())
	}
}
