package main

import (
	"context"
	"io"
	"log"
	"net"
	"strings"

	ordermgt_pb "order-service/server/ecommerce"
	hello_pb "order-service/server/helloworld"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	port           = ":50051"
	orderBatchSize = 3
)

var orderMap = make(map[string]ordermgt_pb.Order)

type orderMgtServer struct {
}

// 단순 RPC
func (s *orderMgtServer) AddOrder(ctx context.Context, orderReq *ordermgt_pb.Order) (*wrappers.StringValue, error) {
	orderMap[orderReq.Id] = *orderReq
	log.Printf("Order : %v -> Added", orderReq.Id)
	return &wrappers.StringValue{Value: "Order Added: " + orderReq.Id}, nil
}

// 단순 RPC
func (s *orderMgtServer) GetOrder(ctx context.Context, orderId *wrappers.StringValue) (*ordermgt_pb.Order, error) {
	ord := orderMap[orderId.Value]
	return &ord, nil
}

// 서버 스트리밍 RPC
func (s *orderMgtServer) SearchOrders(searchQuery *wrappers.StringValue, stream ordermgt_pb.OrderManagement_SearchOrdersServer) error {
	for key, order := range orderMap {
		for _, itemStr := range order.Items {
			if strings.Contains(itemStr, searchQuery.Value) {
				// Send the matching orders in a stream
				log.Printf("Matching Order Found : %s -> Writing Order to the stream ... ", key)
				stream.Send(&order)
				break
			}
		}
	}
	return nil
}

// 서버 스트리밍 RPC
func (s *orderMgtServer) UpdateOrders(stream ordermgt_pb.OrderManagement_UpdateOrdersServer) error {
	ordersStr := "Updated Order IDs : "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			// Finished reading the order stream.
			return stream.SendAndClose(&wrappers.StringValue{Value: "Orders processed " + ordersStr})
		}
		// Update order
		orderMap[order.Id] = *order

		log.Printf("Order ID %v : Updated", order.Id)
		ordersStr += order.Id + ", "
	}
}

// 양방향 스트리밍 RPC
func (s *orderMgtServer) ProcessOrders(stream ordermgt_pb.OrderManagement_ProcessOrdersServer) error {
	batchMarker := 1
	var combinedShipmentMap = make(map[string]ordermgt_pb.CombinedShipment)
	for {
		orderId, err := stream.Recv()
		log.Printf("Reading Proc order ... %v", orderId)
		if err == io.EOF {
			// Client has sent all the messages
			// Send remaining shipments

			log.Printf("EOF %v", orderId)

			for _, comb := range combinedShipmentMap {
				stream.Send(&comb)
			}
			return nil
		}
		if err != nil {
			log.Println(err)
			return err
		}

		destination := orderMap[orderId.GetValue()].Destination
		shipment, found := combinedShipmentMap[destination]

		if found {
			ord := orderMap[orderId.GetValue()]
			shipment.OrdersList = append(shipment.OrdersList, &ord)
			combinedShipmentMap[destination] = shipment
		} else {
			comShip := ordermgt_pb.CombinedShipment{Id: "cmb - " + (orderMap[orderId.GetValue()].Destination), Status: "Processed!"}
			ord := orderMap[orderId.GetValue()]
			comShip.OrdersList = append(shipment.OrdersList, &ord)
			combinedShipmentMap[destination] = comShip
			log.Print(len(comShip.OrdersList), comShip.GetId())
		}

		if batchMarker == orderBatchSize {
			for _, comb := range combinedShipmentMap {
				log.Printf("Shipping : %v -> %d", comb.Id, len(comb.OrdersList))
				stream.Send(&comb)
			}
			batchMarker = 0
			combinedShipmentMap = make(map[string]ordermgt_pb.CombinedShipment)
		} else {
			batchMarker++
		}
	}
}

type helloServer struct {
}

// SayHello implements helloworld.GreeterServer
func (s *helloServer) SayHello(ctx context.Context, in *hello_pb.HelloRequest) (*hello_pb.HelloReply, error) {
	log.Printf("Greeter Service - SayHello RPC")
	return &hello_pb.HelloReply{Message: "Hello " + in.Name}, nil
}

//---------------------------------------------------------
// 코드 5-9 부분
//---------------------------------------------------------
func main() {
	initSampleData()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// gRPC orderMgtServer에 주문 관리 서비스 등록
	ordermgt_pb.RegisterOrderManagementServer(grpcServer, &orderMgtServer{})

	// gRPC orderMgtServer에 Greeter 서비스 등록
	hello_pb.RegisterGreeterServer(grpcServer, &helloServer{})

	//---------------------------------------------------------

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initSampleData() {
	orderMap["102"] = ordermgt_pb.Order{Id: "102", Items: []string{"Google Pixel 3A", "Mac Book Pro"}, Destination: "Mountain View, CA", Price: 1800.00}
	orderMap["103"] = ordermgt_pb.Order{Id: "103", Items: []string{"Apple Watch S4"}, Destination: "San Jose, CA", Price: 400.00}
	orderMap["104"] = ordermgt_pb.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub"}, Destination: "Mountain View, CA", Price: 400.00}
	orderMap["105"] = ordermgt_pb.Order{Id: "105", Items: []string{"Amazon Echo"}, Destination: "San Jose, CA", Price: 30.00}
	orderMap["106"] = ordermgt_pb.Order{Id: "106", Items: []string{"Amazon Echo", "Apple iPhone XS"}, Destination: "Mountain View, CA", Price: 30.00}
}
