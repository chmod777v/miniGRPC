package interfaces

import (
	"context"
	g_serv "grpc/pkg/proto"

	"github.com/jackc/pgx/v5/pgconn"
)

//go:generate mockgen -source=interfaces.go -destination=mock.go
type DatabaseInterface interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) interface{}
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Close()
}

type ServiceServerInterface interface {
	Get(ctx context.Context, req *g_serv.GetRequest) (*g_serv.GetResponse, error)
	Post(ctx context.Context, req *g_serv.PostRequest) (*g_serv.PostResponse, error)
}

type GRPCClientInterface interface {
	GetRequestGrpc(ctx context.Context, requestData *g_serv.GetRequest) (*g_serv.GetResponse, error)
	PostRequestGrpc(ctx context.Context, requestData *g_serv.PostRequest) (*g_serv.PostResponse, error)
	Close()
}

type HTTPServerInterface interface {
	GetRequestGrpc(ctx context.Context, requestData *g_serv.GetRequest) (*g_serv.GetResponse, error)
	PostRequestGrpc(ctx context.Context, requestData *g_serv.PostRequest) (*g_serv.PostResponse, error)
	Close()
}
