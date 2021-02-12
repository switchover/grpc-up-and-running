package main

import (
	"context"
	"io"
	"log"
	"net"

	pb "some-service/server/some"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
}

func (s *server) SomeRPC(ctx context.Context,
	in *pb.SomeRequest) (*pb.SomeResponse, error) {

	log.Print("Requested : ", in.Data)

	return &pb.SomeResponse{Data: "Response"}, nil
}

func (s *server) SomeStreamingRPC(
	stream pb.Service_SomeStreamingRPCServer) error {
	for {
		response, err := stream.Recv()

		if err == io.EOF {
			log.Print("EOF")
			break
		}
		if err == nil {
			log.Print("Search Result : ", response)
		}
	}
	return nil
}

func main() {
	list, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &server{})

	if err := s.Serve(list); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
