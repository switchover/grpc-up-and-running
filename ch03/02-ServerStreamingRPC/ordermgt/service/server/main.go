package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	pb "ordermgt/service/ecommerce"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var orderMap = make(map[string]pb.Order)

type server struct {
}

//---------------------------------------------------------
// 코드 3-5 부분
//---------------------------------------------------------
func (s *server) SearchOrders(searchQuery *wrappers.StringValue,
	stream pb.OrderManagement_SearchOrdersServer) error {
	for key, order := range orderMap {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			if strings.Contains(
				itemStr, searchQuery.Value) {
				// 매칭되는 주문 정보를 스트림에 전송
				err := stream.Send(&order)
				if err != nil {
					return fmt.Errorf(
						"error sending message to stream : %v", err)
				}
				log.Print("Matching Order Found : " + key)
				break
			}
		}
	}
	return nil
}

//---------------------------------------------------------

func main() {
	initSampleData()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderManagementServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initSampleData() {
	orderMap["102"] = pb.Order{Id: "102", Items: []string{"Google Pixel 3A", "Mac Book Pro"}, Destination: "Mountain View, CA", Price: 1800.00}
	orderMap["103"] = pb.Order{Id: "103", Items: []string{"Apple Watch S4"}, Destination: "San Jose, CA", Price: 400.00}
	orderMap["104"] = pb.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub"}, Destination: "Mountain View, CA", Price: 400.00}
	orderMap["105"] = pb.Order{Id: "105", Items: []string{"Amazon Echo"}, Destination: "San Jose, CA", Price: 30.00}
	orderMap["106"] = pb.Order{Id: "106", Items: []string{"Amazon Echo", "Apple iPhone XS"}, Destination: "Mountain View, CA", Price: 300.00}
}
