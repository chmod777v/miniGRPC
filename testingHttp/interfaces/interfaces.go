package httptests

import (
	"context"
	g_serv "grpc/pkg/proto"
)

//go:generate mockgen -source=interfaces.go -destination=mock.go
type HTTPServer interface {
	GetRequestGrpc(ctx context.Context, requestData *g_serv.GetRequest) (*g_serv.GetResponse, error)
	PostRequestGrpc(ctx context.Context, requestData *g_serv.PostRequest) (*g_serv.PostResponse, error)
	Close()
}
