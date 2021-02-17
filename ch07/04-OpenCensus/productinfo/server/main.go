//---------------------------------------------------------
// 코드 7-8 부분
//---------------------------------------------------------
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"

	pb "productinfo/server/ecommerce"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
)

var (
	port = ":50051"
)

// ecommerce/product_info를 구현하는 서버
type server struct {
	productMap map[string]*pb.Product
}

func main() {
	go func() {
		mux := http.NewServeMux()
		zpages.Handle(mux, "/debug")
		log.Fatal(http.ListenAndServe("127.0.0.1:8081", mux))
	}()

	view.RegisterExporter(&exporter.PrintExporter{})

	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	pb.RegisterProductInfoServer(grpcServer, &server{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//---------------------------------------------------------

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

func (s *server) GetProduct(ctx context.Context, in *wrapper.StringValue) (*pb.Product, error) {
	value, exists := s.productMap[in.Value]
	if exists {
		return value, nil
	}
	return nil, errors.New("Product does not exist for the ID" + in.Value)
}
