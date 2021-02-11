package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "order-service/client/ecommerce"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// 서버와의 연결을 구성한다.
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	//---------------------------------------------------------
	// 코드 5-3, 5-4 부분
	//---------------------------------------------------------
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(orderUnaryClientInterceptor), // 5-3 부분
		grpc.WithStreamInterceptor(clientStreamInterceptor))    // 5-4 부분
	//---------------------------------------------------------
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Add Order
	order1 := pb.Order{Id: "101", Items: []string{"iPhone XS", "Mac Book Pro"}, Destination: "San Jose, CA", Price: 2300.00}
	res, _ := c.AddOrder(ctx, &order1)
	log.Print("AddOrder Response -> ", res.Value)

	// Get Order
	retrievedOrder, err := c.GetOrder(ctx, &wrappers.StringValue{Value: "106"})
	log.Print("GetOrder Response -> ", retrievedOrder)

	// Search Order
	searchStream, _ := c.SearchOrders(ctx, &wrappers.StringValue{Value: "Google"})
	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			log.Print("EOF")
			break
		}

		if err == nil {
			log.Print("Search Result : ", searchOrder)
		}
	}

	// Update Orders
	updOrder1 := pb.Order{Id: "102", Items: []string{"Google Pixel 3A", "Google Pixel Book"}, Destination: "Mountain View, CA", Price: 1100.00}
	updOrder2 := pb.Order{Id: "103", Items: []string{"Apple Watch S4", "Mac Book Pro", "iPad Pro"}, Destination: "San Jose, CA", Price: 2800.00}
	updOrder3 := pb.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub", "iPad Mini"}, Destination: "Mountain View, CA", Price: 2200.00}

	updateStream, _ := c.UpdateOrders(ctx)
	_ = updateStream.Send(&updOrder1)
	_ = updateStream.Send(&updOrder2)
	_ = updateStream.Send(&updOrder3)

	updateRes, _ := updateStream.CloseAndRecv()
	log.Printf("Update Orders Res : %v", updateRes)

	// Process Order
	streamProcOrder, _ := c.ProcessOrders(ctx)
	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "102"})
	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "103"})
	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "104"})

	channel := make(chan bool, 1)
	go asynClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "101"})
	_ = streamProcOrder.CloseSend()

	<-channel
}

func asynClientBidirectionalRPC(streamProcOrder pb.OrderManagement_ProcessOrdersClient, c chan bool) {
	for {
		combinedShipment, errProcOrder := streamProcOrder.Recv()
		if errProcOrder == io.EOF {
			break
		}
		log.Printf("Combined shipment : %v", combinedShipment.OrdersList)
	}
	c <- true
}

//---------------------------------------------------------
// 코드 5-3 부분
//---------------------------------------------------------
func orderUnaryClientInterceptor(
	ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 전처리 단계
	log.Println("Method : " + method)

	// 원격 메서드 호출
	err := invoker(ctx, method, req, reply, cc, opts...)

	// 후처리 단계
	log.Println(reply)

	return err
}

//---------------------------------------------------------

//---------------------------------------------------------
// 코드 5-4 부분
//---------------------------------------------------------
func clientStreamInterceptor(
	ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string,
	streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Println("======= [Client Interceptor] ", method)
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStream(s), nil
}

type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("====== [Client Stream Interceptor] "+
		"Receive a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("====== [Client Stream Interceptor] "+
		"Send a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

//---------------------------------------------------------
