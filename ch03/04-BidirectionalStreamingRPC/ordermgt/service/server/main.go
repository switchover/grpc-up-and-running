package main

import (
	"io"
	"log"
	"net"

	pb "ordermgt/service/ecommerce"

	"google.golang.org/grpc"
)

const (
	port           = ":50051"
	orderBatchSize = 3
)

var orderMap = make(map[string]pb.Order)

type server struct {
}

//---------------------------------------------------------
// 코드 3-11 부분
//---------------------------------------------------------
func (s *server) ProcessOrders(
	stream pb.OrderManagement_ProcessOrdersServer) error {
	// ...
	batchMarker := 1
	var combinedShipmentMap = make(map[string]pb.CombinedShipment)

	for {
		orderId, err := stream.Recv()
		if err == io.EOF {
			// ...
			log.Printf("EOF : %s", orderId)
			for _, comb := range combinedShipmentMap {
				stream.Send(&comb)
			}
			return nil
		}
		if err != nil {
			return err
		}
		// 목적지를 기준으로 배송을 구성하는 로직
		// ...
		//---------------------------------------------------------
		destination := orderMap[orderId.GetValue()].Destination
		shipment, found := combinedShipmentMap[destination]

		if found {
			ord := orderMap[orderId.GetValue()]
			shipment.OrdersList = append(shipment.OrdersList, &ord)
			combinedShipmentMap[destination] = shipment
		} else {
			comShip := pb.CombinedShipment{Id: "cmb - " + (orderMap[orderId.GetValue()].Destination), Status: "Processed!"}
			ord := orderMap[orderId.GetValue()]
			comShip.OrdersList = append(shipment.OrdersList, &ord)
			combinedShipmentMap[destination] = comShip
			log.Print(len(comShip.OrdersList), comShip.GetId())
		}
		//---------------------------------------------------------

		//
		if batchMarker == orderBatchSize {
			// 배치 방식으로 클라이언트에게 결합된 주문을 스트리밍한다.
			for _, comb := range combinedShipmentMap {
				// 결합된 배송을 클라이언트에게 전송한다.
				stream.Send(&comb)
			}
			batchMarker = 0
			combinedShipmentMap = make(
				map[string]pb.CombinedShipment)
		} else {
			batchMarker++
		}
	}
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
