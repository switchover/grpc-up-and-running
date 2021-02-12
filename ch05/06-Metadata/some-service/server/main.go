package main

import (
	"context"
	"io"
	"log"
	"net"

	pb "some-service/server/some"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	port = ":50051"
)

type server struct {
}

//---------------------------------------------------------
// 코드 5-13, 5-14 부분
//---------------------------------------------------------
func (s *server) SomeRPC(ctx context.Context,
	in *pb.SomeRequest) (*pb.SomeResponse, error) {
	//-----------------------------------------------------
	// 코드 5-13
	//-----------------------------------------------------
	md, ok := metadata.FromIncomingContext(ctx) // 코드 5-13
	// 메타데이터 활용
	//-----------------------------------------------------
	if ok {
		if t, exists := md["timestamp"]; exists {
			log.Print("timestamp from metadata:")
			for i, e := range t {
				log.Printf("====> Metadata %d : %s", i, e)
			}
		}
	}
	//---------------------------------------------------------
	log.Print("Requested : ", in.Data)

	//-----------------------------------------------------
	// 코드 5-14
	//-----------------------------------------------------
	// 헤더 생성과 전송
	header := metadata.Pairs("header-key", "val")
	grpc.SendHeader(ctx, header)
	// 트레일러 생성과 지정
	trailer := metadata.Pairs("trailer-key", "val")
	grpc.SetTrailer(ctx, trailer)
	//-----------------------------------------------------

	return &pb.SomeResponse{Data: "Response"}, nil
}

//---------------------------------------------------------
// 코드 5-13, 5-14 부분
//---------------------------------------------------------
func (s *server) SomeStreamingRPC(
	stream pb.Service_SomeStreamingRPCServer) error {
	//-----------------------------------------------------
	// 코드 5-13
	//-----------------------------------------------------
	md, ok := metadata.FromIncomingContext(stream.Context())
	// 메타데이터 활용
	//-----------------------------------------------------
	if ok {
		if t, exists := md["timestamp"]; exists {
			log.Print("timestamp from metadata:")
			for i, e := range t {
				log.Printf("====> Metadata %d : %s", i, e)
			}
		}
	}
	//---------------------------------------------------------

	//-----------------------------------------------------
	// 코드 5-14
	//-----------------------------------------------------
	// 헤더 생성과 전송
	header := metadata.Pairs("header-key", "val")
	stream.SendHeader(header)
	// 트레일러 생성과 지정
	trailer := metadata.Pairs("trailer-key", "val")
	stream.SetTrailer(trailer)
	//-----------------------------------------------------

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
