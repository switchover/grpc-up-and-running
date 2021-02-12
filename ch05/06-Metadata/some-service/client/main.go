package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "some-service/client/some"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()

	someRequest := &pb.SomeRequest{Data: "Request"}

	log.Print("Send Header/Trailer Example ...")
	sendMetadata(client, someRequest)

	log.Print("Receive Header/Trailer Example ...")
	receiveMetadata(client, someRequest)
}

func sendMetadata(client pb.ServiceClient, someRequest *pb.SomeRequest) {
	//---------------------------------------------------------
	// 코드 5-11 부분
	//---------------------------------------------------------
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"kn", "vn",
	)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)

	ctxA := metadata.AppendToOutgoingContext(mdCtx,
		"k1", "v1", "k1", "v2", "k2", "v3")

	// 단일 RPC 만들기
	response, err := client.SomeRPC(ctxA, someRequest)
	//---------------------------------------------------------
	if err != nil {
		log.Print("SomeRPC error : ", err.Error())
		os.Exit(-1)
	} else {
		log.Print("Received : ", response.Data)
	}

	//---------------------------------------------------------
	// 코드 5-11 부분
	//---------------------------------------------------------
	// 또는 스트리밍 RPC 만들기
	stream, err := client.SomeStreamingRPC(ctxA)
	//---------------------------------------------------------
	if err != nil {
		log.Print("SomeStreamingRPC error : ", err.Error())
		os.Exit(-1)
	}

	stream.Send(&pb.SomeRequest{Data: "Request 1"})
	stream.Send(&pb.SomeRequest{Data: "Request 2"})

	res, _ := stream.CloseAndRecv()

	log.Print("Res : ", res)
}

func receiveMetadata(client pb.ServiceClient, someRequest *pb.SomeRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//---------------------------------------------------------
	// 코드 5-12 부분
	//---------------------------------------------------------
	var header, trailer metadata.MD

	// ***** 단일 RPC *****

	r, err := client.SomeRPC(
		ctx,
		someRequest,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)

	// 여기서 헤더와 트레일러 맵을 처리한다.
	if t, ok := header["header-key"]; ok {
		log.Print("header-key from header : ", t)
	}
	if t, ok := trailer["trailer-key"]; ok {
		log.Print("trailer-key from trailer : ", t)
	}

	if err != nil {
		log.Print("SomeRPC error : ", err.Error())
		os.Exit(-1)
	} else {
		log.Print("Received : ", r.Data)
	}

	{
		// ***** 스트리밍 RPC *****
		stream, err := client.SomeStreamingRPC(ctx)

		if err != nil {
			log.Print("SomeStreamingRPC error : ", err.Error())
			os.Exit(-1)
		}

		stream.Send(&pb.SomeRequest{Data: "Request 1"})
		stream.Send(&pb.SomeRequest{Data: "Request 2"})

		res, _ := stream.CloseAndRecv()

		log.Print("Res : ", res)

		// 헤더를 조회
		header, err := stream.Header()

		// 트레일러 조회
		trailer := stream.Trailer()

		// 여기서 헤더와 트레일러 맵을 처리한다.
		if t, ok := header["header-key"]; ok {
			log.Print("header-key from header : ", t)
		}
		if t, ok := trailer["trailer-key"]; ok {
			log.Print("trailer-key from trailer : ", t)
		}
		//---------------------------------------------------------
	}
}
