package grpcconect

import (
	"context"
	"fmt"
	g_serv "grpc/pkg/proto"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Server struct {
	grpcClient g_serv.ServiceClient
	conn       *grpc.ClientConn
}
type RequestGrpc interface {
	GetRequestGrpc(ctx context.Context, requestData *g_serv.GetRequest) (*g_serv.GetResponse, error)
	PostRequestGrpc(ctx context.Context, requestData *g_serv.PostRequest) (*g_serv.PostResponse, error)
}

func NewServer(host string, port int) *Server {
	link := fmt.Sprintf("%s:%v", host, port)
	conn, err := grpc.NewClient(link, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to conect to server:", "ERROR", err.Error())
	}
	slog.Info("Conect to gRPC server:", "Host", link)
	return &Server{
		grpcClient: g_serv.NewServiceClient(conn),
		conn:       conn,
	}
}

func (s *Server) GetRequestGrpc(ctx context.Context, requestData *g_serv.GetRequest) (*g_serv.GetResponse, error) {
	response, err := s.grpcClient.Get(ctx, requestData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "gRPC call failed")
	}
	return response, nil
}
func (s *Server) PostRequestGrpc(ctx context.Context, requestData *g_serv.PostRequest) (*g_serv.PostResponse, error) {
	response, err := s.grpcClient.Post(ctx, requestData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "gRPC call failed:%v", err)
	}
	return response, nil
}

func (s *Server) Close() {
	if s.conn != nil {
		slog.Info("ClientConn closing...")
		s.conn.Close()
		slog.Info("ClientConn closed successfully")
	}
}
