# Section 4 : 양방향 스트리밍 RPC (Chapter 3 : gRPC 통신 패턴)

## 예제 코드 리스트
- 코드 3-10 (양방향 스트리밍 서비스 정의 파일) : [order_management.proto](order_management.proto)
- 코드 3-11 (양방향 스트리밍 서비스 구현) : [main.go](ordermgt/service/server/main.go)
- 코드 3-12 (양방향 스트리밍 클라이언트 구현) : [main.go](ordermgt/client/main.go)


---
# 양방향 스트리밍 RPC(Bidirectional Streaming RPC) 서비스 구현 (Go)

## 1. protobuf 정의 파일 수정
[order_management.proto](order_management.proto) (코드 3-10)

## 2. Go 서비스용 모듈 생성
Go 모듈을 위한 디렉토리 생성 후, `go mod` 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
mkdir -p ordermgt/service
cd ordermgt/service
go mod init ordermgt/service
```
※ 실행 위치는 `04-BidirectionalStreaming` 디렉토리를 가정 (이하 동일)

## 3. protobuf 파일 복사
별도로 정의된 `order_management.proto` 파일을 `ecommerce` 디렉토리 생성 후 이 디렉토리로 복사합니다.
```shell
mkdir ecommerce
cp ../../order_management.proto ecommerce
```
- `order_management.proto`는 임의의 위치에서 복사함 (위 예는 현재 예제 디렉토리 구성의 경우임)

## 4. Go 언어 Skeleton 생성 
다음과 같이 이미 설치된 `protoc` 명령을 통해 skeleton 코드를 생성합니다.
```shell
protoc -I ecommerce ecommerce/order_management.proto --go_out=plugins=grpc:ecommerce 
```

## 5. Go 서비스 구현
다음과 같이 Go 서비스를 구현합니다.
[main.go](ordermgt/service/server/main.go) (코드 3-11)
```go
// 일부 코드
func (s *server) ProcessOrders(
	stream pb.OrderManagement_ProcessOrdersServer) error {
	// ...
	for {
		orderId, err := stream.Recv()
		if err == io.EOF {
			// ...
			for _, comb := range combinedShipmentMap {
				stream.Send(&comb)
			}
			return nil
		}
		if err != nil {
			return err
		}
		// 목적지를 기준으로 배송을 구성하는 로직
		// ...
		//
		if batchMarker == orderBatchSize {
			// 배치 방식으로 클라이언트에게 결합된 주문을 스트리밍한다.
			for _, comb := range combinedShipmentMap {
				// 결합된 배송을 클라이언트에게 전송한다.
				stream.Send(&comb)
			}
			batchMarker = 0
			combinedShipmentMap = make(
				map[string]pb.CombinedShipment)
		} else {
			batchMarker++
		}
	}
}
```

## 6. Go 서버 빌드
다음과 같이 서버를 빌드하고 실행합니다.
```shell
go build -i -v -o bin/server server/main.go
```

## 7. Go 클라이언트 생성
다음과 같인 모듈 생성 및 Stub을 생성합니다.
```shell
mkdir -p ordermgt/client
cd ordermgt/client
go mod init ordermgt/client

mkdir ecommerce
cp ../../order_management.proto ecommerce

protoc -I ecommerce ecommerce/order_management.proto --go_out=plugins=grpc:ecommerce 
```

## 8. Go 클라이언트 구현 참조
클라이언트 구현은 다음과 같습니다.
[main.go](ordermgt/client/main.go) (코드 3-12)
```go
	// 일부 코드
	// Process Order
	streamProcOrder, _ := c.ProcessOrders(ctx)
	if err := streamProcOrder.Send(
		&wrapper.StringValue{Value: "102"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "102", err)
	}

	if err := streamProcOrder.Send(
		&wrapper.StringValue{Value: "103"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "103", err)
	}

	if err := streamProcOrder.Send(
		&wrapper.StringValue{Value: "104"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "104", err)
	}

	channel := make(chan struct{})
	go asyncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	if err := streamProcOrder.Send(
		&wrapper.StringValue{Value: "101"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "101", err)
	}

	if err := streamProcOrder.CloseSend(); err != nil {
		log.Fatal(err)
	}

	<-channel
}

func asyncClientBidirectionalRPC(
	streamProcOrder pb.OrderManagement_ProcessOrdersClient,
	c chan struct{}) {
	for {
		combinedShipment, errProcOrder := streamProcOrder.Recv()
		if errProcOrder == io.EOF {
			break
		}
		log.Printf("Combined shipment : ", combinedShipment.OrdersList)
	}
	//<-c
	close(c)
}
```
- 도서 상 예제에는 `func asyncClientBidirectionalRPC()` 함수 정의 부분이 같이 기술되어 있으나, 별도로 구분하여 정리 필요
(예제의 상단 부분은 `main` 함수의 일부분만 표현되어 있음)

## 9. Go 클라이언트 빌드 및 실행
다음과 같이 클라이언트를 빌드 및 실행합니다.
```shell
go build -i -v -o bin/client main.go
bin/client
```


---
# 최종 코드

gRPC의 4가지 통신 패턴에 대한 예제 코드는 원서의 소스 저장소 [3장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch03)을 참고합니다.
