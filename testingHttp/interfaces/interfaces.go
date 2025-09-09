package httptests

import (
	"context"
	g_serv "grpc/pkg/proto"

	"google.golang.org/grpc"
)

//go:generate mockgen -source=interfaces.go -destination=mock.go
type HTTPServer interface {
	Get(ctx context.Context, requestData *g_serv.GetRequest, opts ...grpc.CallOption) (*g_serv.GetResponse, error)
	Post(ctx context.Context, requestData *g_serv.PostRequest, opts ...grpc.CallOption) (*g_serv.PostResponse, error)
	Close()
}
