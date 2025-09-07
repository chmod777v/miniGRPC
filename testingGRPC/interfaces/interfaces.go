package grpctests

import (
	"context"
	"grpc/pkg/database"
	g_serv "grpc/pkg/proto"
)

//go:generate mockgen -source=interfaces.go -destination=mock.go
type Database interface {
	CreatePerson(ctx context.Context, req *g_serv.PostRequest) (int64, error)
	GetPerson(ctx context.Context, id int64) (*database.Person, error)
}

type GRPCServer interface {
	Get(ctx context.Context, req *g_serv.GetRequest) (*g_serv.GetResponse, error)
	Post(ctx context.Context, req *g_serv.PostRequest) (*g_serv.PostResponse, error)
}
