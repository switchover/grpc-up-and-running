package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	pb "productinfo/server/ecommerce"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
)

const (
	address = "localhost:50051"
	bufSize = 1024 * 1024
)

var listener *bufconn.Listener

func initGRPCServerHTTP2() *grpc.Server {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return s
}

//---------------------------------------------------------
// 코드 7-1 부분
//---------------------------------------------------------
func TestServer_AddProduct(t *testing.T) {
	grpcServer := initGRPCServerHTTP2()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		grpcServer.Stop()
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Sumsung S10"
	description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2019"
	price := float32(700.0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx, &pb.Product{Name: name,
		Description: description, Price: price})
	if err != nil {
		t.Fatalf("Could not add product: %v", err)
	}

	if r.Value == "" {
		t.Errorf("Invalid Product ID %s", r.Value)
	}
	log.Printf("Res %s", r.Value)
	grpcServer.Stop()
}

//---------------------------------------------------------
