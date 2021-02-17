//---------------------------------------------------------
// 코드 7-8 부분
//---------------------------------------------------------
package main

import (
	"context"
	"log"
	"time"

	pb "productinfo/client/ecommerce"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	view.RegisterExporter(&exporter.PrintExporter{})

	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(address,
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Can't connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewProductInfoClient(conn)
	// .... // RPC 메서드 호출 생략
	//---------------------------------------------------------

	name := "Sumsung S21"
	description := "Samsung Galaxy S10 is the latest smart phone, launched in Jan. 2021"
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
	log.Printf("Product: %s", product.String())
}
