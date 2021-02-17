package main

import (
	"context"
	"log"
	"time"

	pb "productinfo/ingress/ecommerce"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	address = "productinfo:80"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewProductInfoClient(conn)

	// Contact the server and print out its response.
	name := "Sumsung S21"
	description := "Samsung Galaxy S21 is the latest smart phone, launched in Jan. 2021"
	price := float32(700.0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	product, err := c.GetProduct(ctx, &wrapper.StringValue{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %v", product.String())
}
