package my_grpc

import (
	"context"
	db "grpc/pkg/database"
	g_serv "grpc/pkg/proto"
	"log/slog"
)

type Person struct {
	Id      int
	User_id int
	Name    string
	Admin   bool
}
type Server struct {
	g_serv.UnimplementedServiceServer
	Db *db.Database
}

func (s *Server) Get(ctx context.Context, req *g_serv.GetRequest) (*g_serv.GetResponse, error) {
	slog.Debug("Request", "Method", "Get", "Data", req)

	var person Person
	err := s.Db.Pool.QueryRow(context.Background(), "SELECT * FROM people WHERE ID=$1", req.Id).
		Scan(&person.Id, &person.Name, &person.Admin, &person.User_id)
	if err != nil {
		slog.Error("Error while receiving data", "ERROR", err)
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

	var person Person
	err := s.Db.Pool.QueryRow(context.Background(),
		"INSERT INTO people (User_id, Name, Admin) VALUES ($1, $2, $3) RETURNING id",
		req.Info.UserId, req.Info.Name, req.Info.Admin).Scan(&person.Id)
	if err != nil {
		slog.Error("Error adding field", "ERROR", err)
		return nil, err
	}

	resp := &g_serv.PostResponse{
		Id: int64(person.Id),
		Info: &g_serv.UserInfo{
			Name:   req.Info.Name,
			UserId: req.Info.UserId,
			Admin:  req.Info.Admin,
		},
	}
	return resp, nil
}
