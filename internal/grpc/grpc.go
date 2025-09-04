package my_grpc

import (
	"context"
	db "grpc/pkg/database"
	g_serv "grpc/pkg/proto"
	"log"
	"log/slog"
)

type Server struct {
	g_serv.UnimplementedServiceServer
	Db *db.Database
}

func (s *Server) Get(ctx context.Context, req *g_serv.GetRequest) (*g_serv.GetResponse, error) {
	slog.Debug("Request", "Data", req)
	resp := &g_serv.GetResponse{
		Info: &g_serv.UserInfo{
			Id:      req.GetId(),
			Name:    "Dima",
			IsHuman: true,
		},
	}
	_, err := s.Db.Pool.Exec(context.Background(), "INSERT INTO people (Name, Admin) VALUES ($1, $2)", "Алекс", true)
	if err != nil {
		log.Fatal("Ошибка при вставке данных:", err)
		return nil, err
	}
	return resp, nil

}
