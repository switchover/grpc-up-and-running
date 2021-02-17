package main

import (
	"context"
	"log"
	"net"

	pb "productinfo/server/ecommerce"
	"productinfo/server/tracer"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = ":50051"
)

type server struct {
	productMap map[string]*pb.Product
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterProductInfoServer(grpcServer, &server{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// initialize opencensus jaeger exporter
	tracer.InitTracing()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) AddProduct(ctx context.Context, in *pb.Product) (*wrapper.StringValue, error) {
	out, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[in.Id] = in
	return &wrapper.StringValue{Value: in.Id}, nil
}

//---------------------------------------------------------
// 코드 7-14
//---------------------------------------------------------
// GetProduct는 ecommerce.GetProduct를 구현한다.
func (s *server) GetProduct(ctx context.Context, in *wrapper.StringValue) (
	*pb.Product, error) {
	ctx, span := trace.StartSpan(ctx, "ecommerce.GetProduct")
	defer span.End()
	value, exists := s.productMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}

//---------------------------------------------------------
