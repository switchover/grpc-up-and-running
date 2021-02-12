package main

import (
	"context"
	"log"
	"time"

	pb "some-service/client/some"

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

	client := pb.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	someRequest := &pb.SomeRequest{Data: "Request"}
	response, err := client.SomeRPC(ctx, someRequest)
	if err != nil {
		log.Print("SomeRPC error : ", err.Error())
		return
	} else {
		log.Print("Received : ", response.Data)
	}

	stream, err := client.SomeStreamingRPC(ctx)
	if err != nil {
		log.Print("SomeStreamingRPC error : ", err.Error())
		return
	}

	stream.Send(&pb.SomeRequest{Data: "Request 1"})
	stream.Send(&pb.SomeRequest{Data: "Request 2"})

	res, _ := stream.CloseAndRecv()

	log.Print("Res : ", res)

}
