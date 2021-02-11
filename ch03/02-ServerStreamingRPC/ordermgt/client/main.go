package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "ordermgt/client/ecommerce"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
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

	//---------------------------------------------------------
	// 코드 3-6 부분
	//---------------------------------------------------------
	c := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	searchStream, _ := c.SearchOrders(ctx,
		&wrapper.StringValue{Value: "Google"})

	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			break
		}
		// 기타 가능한 에러의 처리
		log.Print("Search Result : ", searchOrder)
	}
	//---------------------------------------------------------
}
