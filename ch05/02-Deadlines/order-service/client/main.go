package main

import (
	"context"
	"log"
	"time"

	pb "order-service/client/ecommerce"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	address = "localhost:50051"
)

func main() {
	// 서버와의 연결을 구성한다.
	//---------------------------------------------------------
	// 코드 5-5 부분
	//---------------------------------------------------------
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewOrderManagementClient(conn)

	clientDeadline := time.Now().Add(
		time.Duration(2 * time.Second))
	ctx, cancel := context.WithDeadline(
		context.Background(), clientDeadline)

	defer cancel()

	// Order 등록
	order1 := pb.Order{Id: "101",
		Items:       []string{"iPhone XS", "Mac Book Pro"},
		Destination: "San Jose, CA",
		Price:       2300.00}
	res, addErr := client.AddOrder(ctx, &order1)
	if addErr != nil {
		got := status.Code(addErr)
		log.Printf("Error Occured -> addOrder : , %v:", got)
	} else {
		log.Print("AddOrder Response -> ", res.Value)
	}
	//---------------------------------------------------------
}
