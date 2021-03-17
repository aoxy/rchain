package main

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"rchain/examples/grpc/proto"
)

const PORT = "9002"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatal("grpc.Dial err:", err)
	}
	defer conn.Close()

	client := proto.NewStreamServiceClient(conn)

	err = printLists(client, &proto.StreamRequest{Pt: &proto.StreamPoint{
		Name:  "gRPC stream client: List",
		Value: 2020,
	}})
	if err != nil {
		log.Fatal("printLists err:", err)
	}

	err = printRecord(client, &proto.StreamRequest{Pt: &proto.StreamPoint{
		Name:  "gRPC stream client: Record",
		Value: 2020,
	}})
	if err != nil {
		log.Fatal("printRecord err:", err)
	}

	err = printRoute(client, &proto.StreamRequest{Pt: &proto.StreamPoint{
		Name:  "gRPC stream client: Route",
		Value: 2020,
	}})
	if err != nil {
		log.Fatal("printRoute err:", err)
	}
}

func printLists(client proto.StreamServiceClient, r *proto.StreamRequest) error {
	stream, err := client.List(context.TODO(), r)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: pt.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	return nil
}

func printRecord(client proto.StreamServiceClient, r *proto.StreamRequest) error {

	stream, err := client.Record(context.TODO())
	if err != nil {
		return err
	}

	for i := 0; i <= 6; i++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("resp: pt.name: %s, pt.value: %d \n", resp.Pt.Name, resp.Pt.Value)
	return nil
}

func printRoute(client proto.StreamServiceClient, r *proto.StreamRequest) error {

	stream, err := client.Route(context.TODO())
	if err != nil {
		return err
	}

	for i := 0; i <= 6; i++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: pt.Name: %s, pt.Value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	stream.CloseSend()
	return nil
}
