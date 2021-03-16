package main

import (
	"context"
	"github.com/dawn1806/rchain_fork/grpc/example/proto"
	"google.golang.org/grpc"
	"log"
)

const PORT = "9001"

func main() {

	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatal("grpc.Dial err:", err)
	}
	defer conn.Close()

	client := proto.NewSearchServiceClient(conn)
	res, err := client.Search(context.TODO(), &proto.SearchRequest{Request: "gRPC"})
	if err != nil {
		log.Fatal("client.Search err:", err)
	}

	log.Println("res:", res.GetResponse())
}
