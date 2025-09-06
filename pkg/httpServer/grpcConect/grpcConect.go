package grpcconect

import (
	"context"
	"fmt"
	g_serv "grpc/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	grpcClient g_serv.ServiceClient
	conn       *grpc.ClientConn
}

func NewServer() *Server {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Failed to conect to server:" + err.Error())
	}

	return &Server{
		grpcClient: g_serv.NewServiceClient(conn),
		conn:       conn,
	}
}

func (s *Server) GetRequestGrpc(ctx context.Context, requestData *g_serv.GetRequest) (*g_serv.GetResponse, error) {
	response, err := s.grpcClient.Get(ctx, requestData)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", err)
	}
	return response, nil
}
func (s *Server) PostRequestGrpc(ctx context.Context, requestData *g_serv.PostRequest) (*g_serv.PostResponse, error) {
	response, err := s.grpcClient.Post(ctx, requestData)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", err)
	}
	return response, nil
}

func (s *Server) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
