//---------------------------------------------------------
// 코드 6-8 부분
//---------------------------------------------------------
package main

import (
	"context"
	"log"
	"time"

	pb "productinfo/client/ecommerce"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	crtFile  = "server.crt"
)

func main() {
	auth := oauth.NewOauthAccess(fetchToken())

	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
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

	// .... // RPC 메소드 호출을 건너뛴다.
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
// 코드 6-8 부분
//---------------------------------------------------------
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}

//---------------------------------------------------------
