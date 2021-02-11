package main

import (
	"context"
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
	// 코드 3-3 부분
	//---------------------------------------------------------
	orderMgtClient := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 주문 정보 가져오기
	retrievedOrder, rr := orderMgtClient.GetOrder(ctx,
		&wrapper.StringValue{Value: "106"})
	log.Print("GetOrder Response -> : ", retrievedOrder)
	//---------------------------------------------------------

	if rr != nil {
		log.Fatalf("GetOrder() error : %v", rr)
	}
}
