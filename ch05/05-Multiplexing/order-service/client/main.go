package main

import (
	"context"
	"log"
	"time"

	pb "order-service/client/ecommerce"
	hwpb "order-service/client/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	address = "localhost:50051"
)

func main() {
	//---------------------------------------------------------
	// 코드 5-10 부분
	//---------------------------------------------------------
	// 서버에 대한 연결을 설정한다.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderManagementClient := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Add Order
	order1 := pb.Order{Id: "101", Items: []string{"iPhone XS", "Mac Book Pro"}, Destination: "San Jose, CA", Price: 2300.00}
	res, addErr := orderManagementClient.AddOrder(ctx, &order1)

	if addErr != nil {
		got := status.Code(addErr)
		log.Printf("Error Occured -> addOrder : %v", got)
	} else {
		log.Printf("AddOrder Response -> %v", res.Value)
	}

	helloClient := hwpb.NewGreeterClient(conn)

	hwcCtx, hwcCancel := context.WithTimeout(context.Background(), time.Second)
	defer hwcCancel()

	// Say hello RPC
	helloResponse, err := helloClient.SayHello(hwcCtx,
		&hwpb.HelloRequest{Name: "gRPC Up and Running!"})

	log.Print("Greeting: ", helloResponse.Message)
	//---------------------------------------------------------

}
