package main

import (
	"context"
	"github.com/dawn1806/rchain_fork/grpc/example/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *proto.SearchRequest) (*proto.SearchResponse, error) {
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
