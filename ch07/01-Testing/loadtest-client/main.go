package main

import (
	"context"
	"log"
	"time"

	pb "loadtest-client/helloworld"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Vincent"})
	if err != nil {
		log.Fatalf("Could not say hello: %v", err)
	}
	log.Printf("Reply message: %v", r.Message)
}
