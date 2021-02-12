package main

import (
	"context"
	"log"
	"time"

	pb "order-service/client/ecommerce"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
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

	c := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Add Order
	order1 := pb.Order{Id: "101", Items: []string{"iPhone XS", "Mac Book Pro"}, Destination: "San Jose, CA", Price: 2300.00}
	res, _ := c.AddOrder(ctx, &order1, grpc.UseCompressor(gzip.Name))
	log.Print("AddOrder Response -> ", res.Value)
}
