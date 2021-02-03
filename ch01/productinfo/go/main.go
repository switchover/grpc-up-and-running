package main

// 코드 1-2 부분
import (
	"context"
	"errors"
	"log"
	"net"

	pb "github.com/grpc-up-and-running/samples/ch02/productinfo/go/proto"
	"google.golang.org/grpc"
)

type server struct {
	productMap map[string]*pb.Product
}

// Go 언어를 사용한 ProductInfo 구현

// 제품 등록을 위한 원격 메서드
func (s *server) AddProduct(ctx context.Context, in *pb.Product) (*pb.ProductID, error) {
	// 업무 로직
	in.Id = "Product1"

	log.Println("Name :", in.Name)
	log.Println("Desc. :", in.Description)

	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[in.Id] = in

	return &pb.ProductID{Value: in.Id}, nil
}

// 제품 조회용 원격 메서드
func (s *server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.Product, error) {
	// 업무 로직

	product, exists := s.productMap[in.Value]
	if exists && product != nil {
		return product, nil
	}
	return nil, errors.New("product not found")
}

// 서비스 port
const port = ":8080"

// 코드 1-3 부분
func main() {
	lis, _ := net.Listen("tcp", port)
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
