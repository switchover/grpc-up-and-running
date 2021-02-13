//---------------------------------------------------------
// 코드 6-6 부분
//---------------------------------------------------------
package main

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	pb "productinfo/client/ecommerce"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	crtFile  = "server.crt"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	auth := basicAuth{
		username: "admin",
		password: "admin",
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(auth),
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewProductInfoClient(conn)

	// .... // RPC 메서드 호출 생략
	//---------------------------------------------------------

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

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %s", product.String())
}

//---------------------------------------------------------
// 코드 6-5 부분
//---------------------------------------------------------
type basicAuth struct {
	username string
	password string
}

func (b basicAuth) GetRequestMetadata(ctx context.Context,
	in ...string) (map[string]string, error) {
	auth := b.username + ":" + b.password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

func (b basicAuth) RequireTransportSecurity() bool {
	return true
}

//---------------------------------------------------------
