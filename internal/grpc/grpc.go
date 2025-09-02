package my_grpc

import (
	"context"
	g_serv "grpc/pkg/proto"
	"log/slog"
)

type Server struct {
	g_serv.UnimplementedServiceServer
}

func (s *Server) Get(ctx context.Context, req *g_serv.GetRequest) (*g_serv.GetResponse, error) {
	slog.Debug("Request", "Id", req.GetId())
	resp := &g_serv.GetResponse{
		Info: &g_serv.UserInfo{
			Id:      req.GetId(),
			Name:    "Dima",
			IsHuman: true,
		},
	}
	return resp, nil

}
