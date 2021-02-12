package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "order-service/client/ecommerce"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// 서버와의 연결을 구성한다.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)
	//---------------------------------------------------------
	// 코드 5-6 부분
	//---------------------------------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	streamProcOrder, _ := client.ProcessOrders(ctx)
	_ = streamProcOrder.Send(&wrapper.StringValue{Value: "102"})
	_ = streamProcOrder.Send(&wrapper.StringValue{Value: "103"})
	_ = streamProcOrder.Send(&wrapper.StringValue{Value: "104"})

	channel := make(chan bool, 1)
	go asncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	// RPC 취소
	cancel()
	log.Printf("RPC Status : %s", ctx.Err())

	_ = streamProcOrder.Send(&wrapper.StringValue{Value: "101"})
	_ = streamProcOrder.CloseSend()

	<-channel
	//---------------------------------------------------------
}

//---------------------------------------------------------
// 코드 5-6 부분
//---------------------------------------------------------
func asncClientBidirectionalRPC(
	streamProcOrder pb.OrderManagement_ProcessOrdersClient, c chan bool) {
	for {
		combinedShipment, errProcOrder := streamProcOrder.Recv()
		if errProcOrder != nil {
			log.Printf("Error Receiving messages %v", errProcOrder)
			break
		} else {
			if errProcOrder == io.EOF {
				break
			}
			log.Printf("Combined shipment : %s", combinedShipment.OrdersList)
		}
	}
	c <- true
}

//---------------------------------------------------------
