// 코드 2-6 부분
package main

import (
	"context"

	pb "productinfo/service/ecommerce"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server는 ecommerce/product_info를 구현하는 데 사용된다.
type server struct {
	productMap map[string]*pb.Product
}

// AddProduct는 ecommerce.AddProduct를 구현한다.
func (s *server) AddProduct(ctx context.Context,
	in *pb.Product) (*pb.ProductID, error) {
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"Error while generating Product ID", err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[in.Id] = in
	return &pb.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

// GetProduct는 ecommerce.GetProduct를 구현한다.
func (s *server) GetProduct(ctx context.Context,
	in *pb.ProductID) (*pb.Product, error) {

	value, exists := s.productMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}
