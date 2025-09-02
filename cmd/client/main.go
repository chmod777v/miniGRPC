package main

import (
	"context"
	g_serv "grpc/pkg/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Failed to conect to server:" + err.Error())
	}
	defer conn.Close()

	cli := g_serv.NewServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, &g_serv.GetRequest{Id: 123})
	if err != nil {
		log.Fatalln("failed to get response:", err.Error())
	}
	log.Println(resp.GetInfo())
}
