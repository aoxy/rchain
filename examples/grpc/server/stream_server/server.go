package main

import (
	"github.com/dawn1806/rchain_fork/grpc/example/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type StreamService struct{}

const PORT = "9002"

func main() {
	server := grpc.NewServer()
	proto.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatal("net.Listen err:", err)
	}

	server.Serve(lis)
}

func (s *StreamService) List(r *proto.StreamRequest, stream proto.StreamService_ListServer) error {

	for i := 0; i <= 6; i++ {
		err := stream.Send(&proto.StreamResponse{
			Pt: &proto.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(i),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StreamService) Record(stream proto.StreamService_RecordServer) error {

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.StreamResponse{Pt: &proto.StreamPoint{
				Name:  "gRPC stream Server: Record",
				Value: 1,
			}})
		}
		if err != nil {
			return err
		}

		log.Printf("stream.Recv: name: %s, value: %d \n", r.Pt.Name, r.Pt.Value)
	}
}

func (s *StreamService) Route(stream proto.StreamService_RouteServer) error {

	n := 0
	for {
		err := stream.Send(&proto.StreamResponse{Pt: &proto.StreamPoint{
			Name:  "gRPC stream server: Route",
			Value: int32(n),
		}})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++

		log.Printf("stream.Recv: pt.Name: %s, pt.Value: %d \n", r.Pt.Name, r.Pt.Value)
	}

}
