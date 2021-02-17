//---------------------------------------------------------
// 코드 7-15 부분
//---------------------------------------------------------
package main

import (
	"context"
	"log"

	pb "productinfo/client/ecommerce"
	"productinfo/client/tracer"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	tracer.InitTracing()

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	ctx, span := trace.StartSpan(context.Background(),
		"ecommerce.ProductInfoClient")

	name := "Apple iphone 11"
	description := "Apple iphone 11 is the latest smartphone, launched in September 2019"
	price := float32(700.0)
	r, err := c.AddProduct(ctx, &pb.Product{Name: name,
		Description: description, Price: price})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	//product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	product, err := c.GetProduct(ctx, &wrapper.StringValue{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %v", product.String())
	span.End()
}

//---------------------------------------------------------
