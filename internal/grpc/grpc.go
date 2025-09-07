package my_grpc

import (
	"context"
	"grpc/pkg/database"
	g_serv "grpc/pkg/proto"
	"log/slog"
)

type Server struct {
	g_serv.UnimplementedServiceServer
	Db database.DatabaseCG
}

func (s *Server) Get(ctx context.Context, req *g_serv.GetRequest) (*g_serv.GetResponse, error) {
	slog.Debug("Request", "Method", "Get", "Data", req)
	person, err := s.Db.GetPerson(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	resp := &g_serv.GetResponse{
		Info: &g_serv.UserInfo{
			Name:   person.Name,
			UserId: int64(person.User_id),
			Admin:  person.Admin,
		},
	}
	return resp, nil
}
func (s *Server) Post(ctx context.Context, req *g_serv.PostRequest) (*g_serv.PostResponse, error) {
	slog.Debug("Request", "Method", "Post", "Data", req)
	id, err := s.Db.CreatePerson(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := &g_serv.PostResponse{
		Id: id,
		Info: &g_serv.UserInfo{
			Name:   req.Info.Name,
			UserId: req.Info.UserId,
			Admin:  req.Info.Admin,
		},
	}
	return resp, nil
}
