package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"rchain/examples/grpc/proto"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *proto.SearchRequest) (*proto.SearchResponse, error) {
	fmt.Println("接收到请求了")
	return &proto.SearchResponse{
		Response: r.GetRequest() + " Server",
	}, nil
}

const PORT = "9001"

func main() {
	server := grpc.NewServer()
	proto.RegisterSearchServiceServer(server, &SearchService{})

	conn, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatal("net.Listen err:", err)
	}

	server.Serve(conn)
}
