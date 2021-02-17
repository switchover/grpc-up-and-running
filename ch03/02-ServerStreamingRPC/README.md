# Section 2 : 서버 스트리밍 RPC (Chapter 3 : gRPC 통신 패턴)

## 예제 코드 리스트
- 코드 3-4 (서버 스트리밍 서비스 정의 파일) : [order_management.proto](order_management.proto)
- 코드 3-5 (서버 스트리밍 서비스 구현) : [main.go](ordermgt/service/server/main.go)
- 코드 3-6 (서버 스트리밍 클라이언트 구현) : [main.go](ordermgt/client/main.go)


---
# 서버 스트리밍 RPC(Server Streaming RPC) 서비스 구현 (Go)

## 1. protobuf 정의 파일 수정
[order_management.proto](order_management.proto) (코드 3-4)

## 2. Go 서비스용 모듈 생성
Go 모듈을 위한 디렉토리 생성 후, `go mod` 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
$ mkdir -p ordermgt/service
$ cd ordermgt/service
$ go mod init ordermgt/service
```
※ 실행 위치는 `02-ServerStreaming` 디렉토리를 가정 (이하 동일)

## 3. protobuf 파일 복사
별도로 정의된 `order_management.proto` 파일을 `ecommerce` 디렉토리 생성 후 이 디렉토리로 복사합니다.
```shell
$ mkdir ecommerce
$ cp ../../order_management.proto ecommerce
```
- `order_management.proto`는 임의의 위치에서 복사함 (위 예는 현재 예제 디렉토리 구성의 경우임)

## 4. Go 언어 Skeleton 생성 
다음과 같이 이미 설치된 `protoc` 명령을 통해 skeleton 코드를 생성합니다.
```shell
$ protoc -I ecommerce ecommerce/order_management.proto --go_out=plugins=grpc:ecommerce 
```

## 5. Go 서비스 구현
다음과 같이 Go 서비스를 구현합니다.
[main.go](ordermgt/service/server/main.go) (코드 3-5)
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
$ go build -i -v -o bin/server server/main.go
```

## 7. Go 클라이언트 생성
다음과 같인 모듈 생성 및 Stub을 생성합니다.
```shell
$ mkdir -p ordermgt/client
$ cd ordermgt/client
$ go mod init ordermgt/client

$ mkdir ecommerce
$ cp ../../order_management.proto ecommerce

$ protoc -I ecommerce ecommerce/order_management.proto --go_out=plugins=grpc:ecommerce 
```

## 8. Go 클라이언트 구현 참조
클라이언트 구현은 다음과 같습니다.
[main.go](ordermgt/client/main.go) (코드 3-6)
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

## 9. Go 클라이언트 빌드 및 실행
다음과 같이 클라이언트를 빌드 및 실행합니다.
```shell
$ go build -i -v -o bin/client main.go
$ bin/client
```
