# Chapter 3 : gRPC 통신 패턴

## 예제 코드 리스트
- 코드 3-1 : [order_management.proto](01-SimpleRPC/order_management.proto)
- 코드 3-2 : [main.go](01-SimpleRPC/ordermgt/service/server/main.go)
- 코드 3-3 : [main.go](01-SimpleRPC/ordermgt/client/main.go)
- 코드 3-4 : [order_management.proto](02-ServerStreamingRPC/order_management.proto)
- 코드 3-5 : [main.go](02-ServerStreamingRPC/ordermgt/service/server/main.go)
- 코드 3-6 : [main.go](02-ServerStreamingRPC/ordermgt/client/main.go)
- 코드 3-7 : [order_management.proto](03-ClientStreamingRPC/order_management.proto)
- 코드 3-8 : [main.go](03-ClientStreamingRPC/ordermgt/service/server/main.go)
- 코드 3-9 : [main.go](03-ClientStreamingRPC/ordermgt/client/main.go)
- 코드 3-10 : [order_management.proto](04-BidirectionalStreamingRPC/order_management.proto)
- 코드 3-11 : [main.go](04-BidirectionalStreamingRPC/ordermgt/service/server/main.go)
- 코드 3-12 : [main.go](04-BidirectionalStreamingRPC/ordermgt/client/main.go)

## 정오
### 코드 부분
- 98 페이지 코드 3-9. `client` 정의 부분 :
    ```
    c := pb.NewOrderManagementClient(conn)
    ```
    :arrow_right:
    ```
    client := pb.NewOrderManagementClient(conn)
    ```

---
# 단순 RPC(Simple RPC) 서비스 구현 (Go)

## 1. protobuf 정의 파일 생성
[order_management.proto](01-SimpleRPC/order_management.proto) (코드 3-1)

## 2. Go 서비스용 모듈 생성
Go 모듈을 위한 디렉토리 생성 후, `go mod` 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
mkdir -p ordermgt/service
cd ordermgt/service
go mod init ordermgt/service
```
※ 실행 위치는 `01-SimpleRPC` 디렉토리를 가정 (이하 동일)

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
[main.go](01-SimpleRPC/ordermgt/service/server/main.go) (코드 3-2)
```go
// 일부 코드
func (s *server) GetOrder(ctx context.Context,
	orderId *wrapper.StringValue) (*pb.Order, error) {
	// 서비스 구현
	ord := orderMap[orderId.Value]
	return &ord, nil
}
```

## 6. Go 서버 빌드
다음과 같이 서버를 빌드하고 실행합니다.
```shell
go build -i -v -o bin/server server/main.go
```

## 7. Go 클라이언트 구현 참조
클라이언트 구현은 다음과 같습니다.
[main.go](01-SimpleRPC/ordermgt/client/main.go) (코드 3-3)
```go
// 일부 코드
// 서버와의 연결을 구성한다.
// ...
orderMgtClient := pb.NewOrderManagementClient(conn)
// ...

// 주문 정보 가져오기
retrievedOrder, rr := orderMgtClient.GetOrder(ctx,
    &wrapper.StringValue{Value: "106"})
log.Print("GetOrder Response -> : ", retrievedOrder)
// ...
```

## 8. Go 클라이언트 빌드 및 실행
다음과 같이 클라이언트를 빌드 및 실행합니다.
```shell
go build -i -v -o bin/client main.go
bin/client
```

---
# 서버 스트리밍 RPC(Server Streaming RPC) 서비스 구현 (Go)

## 1. protobuf 정의 파일 수정
[order_management.proto](02-ServerStreamingRPC/order_management.proto) (코드 3-4)

## 2. Go 서비스용 모듈 생성
Go 모듈을 위한 디렉토리 생성 후, `go mod` 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
mkdir -p ordermgt/service
cd ordermgt/service
go mod init ordermgt/service
```
※ 실행 위치는 `02-ServerStreaming` 디렉토리를 가정 (이하 동일)

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
[main.go](02-ServerStreamingRPC/ordermgt/service/server/main.go) (코드 3-5)
```go
// 일부 코드
func (s *server) SearchOrders(searchQuery *wrappers.StringValue,
	stream pb.OrderManagement_SearchOrdersServer) error {
	for key, order := range orderMap {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			if strings.Contains(
				itemStr, searchQuery.Value) {
				// 매칭되는 주문 정보를 스트림에 전송
				err := stream.Send(&order)
				if err != nil {
					return fmt.Errorf(
						"error sending message to stream : %v", err)
				}
				log.Print("Matching Order Found : " + key)
				break
			}
		}
	}
	return nil
}
```

## 6. Go 서버 빌드
다음과 같이 서버를 빌드하고 실행합니다.
```shell
go build -i -v -o bin/server server/main.go
```

## 7. Go 클라이언트 구현 참조
클라이언트 구현은 다음과 같습니다.
[main.go](02-ServerStreamingRPC/ordermgt/client/main.go) (코드 3-6)
```go
// 일부 코드
// 서버와의 연결을 구성한다.
// ...
c := pb.NewOrderManagementClient(conn)
// ...

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
```

## 8. Go 클라이언트 빌드 및 실행
다음과 같이 클라이언트를 빌드 및 실행합니다.
```shell
go build -i -v -o bin/client main.go
bin/client
```


# 클라이언트 스트리밍 RPC(Client Streaming RPC) 서비스 구현 (Go)

## 1. protobuf 정의 파일 수정
[order_management.proto](03-ClientStreamingRPC/order_management.proto) (코드 3-7)

## 2. Go 서비스용 모듈 생성
Go 모듈을 위한 디렉토리 생성 후, `go mod` 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
mkdir -p ordermgt/service
cd ordermgt/service
go mod init ordermgt/service
```
※ 실행 위치는 `03-ClientStreaming` 디렉토리를 가정 (이하 동일)

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
[main.go](03-ClientStreamingRPC/ordermgt/service/server/main.go) (코드 3-8)
```go
// 일부 코드
func (s *server) UpdateOrders(stream pb.OrderManagement_UpdateOrdersServer) error {
	ordersStr := "Updated Order IDs : "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			// 주문 스트림 읽기를 종료한다.
			return stream.SendAndClose(
				&wrapper.StringValue{Value: "Orders processed " + ordersStr})
		}
		// Update order
		orderMap[order.Id] = *order
		log.Printf("Order ID ", order.Id, ": Updated")
		ordersStr += order.Id + ", "
	}
}
```

## 6. Go 서버 빌드
다음과 같이 서버를 빌드하고 실행합니다.
```shell
go build -i -v -o bin/server server/main.go
```

## 7. Go 클라이언트 구현 참조
클라이언트 구현은 다음과 같습니다.
[main.go](03-ClientStreamingRPC/ordermgt/client/main.go) (코드 3-9)
```go
	// 일부 코드
	// 서버와의 연결을 구성한다.
	// ...
	//c := pb.NewOrderManagementClient(conn)
	client := pb.NewOrderManagementClient(conn)
	// ...
	
	updateStream, err := client.UpdateOrders(ctx)

	if err != nil {
		log.Fatalf("%v.UpdateOrders(_) = _, %v", client, err)
	}

	// Updating order 1
	if err := updateStream.Send(&updOrder1); err != nil {
		log.Fatalf("%v.Send(%v) = %v",
			updateStream, updOrder1, err)
	}

	// Updating order 2
	if err := updateStream.Send(&updOrder2); err != nil {
		log.Fatalf("%v.Send(%v) = %v",
			updateStream, updOrder2, err)
	}

	// Updating order 3
	if err := updateStream.Send(&updOrder3); err != nil {
		log.Fatalf("%v.Send(%v) = %v",
			updateStream, updOrder3, err)
	}

	updateRes, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v",
			updateStream, err, nil)
	}
	log.Printf("Update Orders Res : %s", updateRes)
```

## 8. Go 클라이언트 빌드 및 실행
다음과 같이 클라이언트를 빌드 및 실행합니다.
```shell
go build -i -v -o bin/client main.go
bin/client
```


# 양방향 스트리밍 RPC(Bidirectional Streaming RPC) 서비스 구현 (Go)

## 1. protobuf 정의 파일 수정
[order_management.proto](04-BidirectionalStreamingRPC/order_management.proto) (코드 3-10)

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
[main.go](04-BidirectionalStreamingRPC/ordermgt/service/server/main.go) (코드 3-11)
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

## 7. Go 클라이언트 구현 참조
클라이언트 구현은 다음과 같습니다.
[main.go](04-BidirectionalStreamingRPC/ordermgt/client/main.go) (코드 3-12)
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
	<-c
}
```

## 8. Go 클라이언트 빌드 및 실행
다음과 같이 클라이언트를 빌드 및 실행합니다.
```shell
go build -i -v -o bin/client main.go
bin/client
```


---
# 최종 코드

gRPC의 4가지 통신 패턴에 대한 예제 코드는 원서의 소스 저장소 [3장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch03)을 참고합니다.
