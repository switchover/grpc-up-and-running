//---------------------------------------------------------
// 코드 7-10 부분
//---------------------------------------------------------
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "productinfo/server/ecommerce"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	reg = prometheus.NewRegistry()

	grpcMetrics = grpc_prometheus.NewServerMetrics()

	customMetricCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "product_mgt_server_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
)

func init() {
	reg.MustRegister(grpcMetrics, customMetricCounter)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9092)}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	pb.RegisterProductInfoServer(grpcServer, &server{})
	grpcMetrics.InitializeMetrics(grpcServer)

	// 프로메테우스에 대한 http 서버를 시작한다.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//---------------------------------------------------------

var (
	port = ":50051"
)

// ecommerce/product_info를 구현하는 서버
type server struct {
	productMap map[string]*pb.Product
}

//---------------------------------------------------------
// 코드 7-11 부분
//---------------------------------------------------------
// AddProduct는 ecommerce.AddProduct를 구현한다.
func (s *server) AddProduct(ctx context.Context,
	in *pb.Product) (*wrapper.StringValue, error) {
	customMetricCounter.WithLabelValues(in.Name).Inc()
	//---------------------------------------------------------
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
