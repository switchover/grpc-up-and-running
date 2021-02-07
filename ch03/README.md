# Chapter 3 : gRPC 통신 패턴

## 예제 코드 리스트
- 코드 3-1 : [order_management.proto](01-SimpleRPC/order_management.proto)
- 코드 3-2 : [main.go](01-SimpleRPC/ordermgt/service/server/main.go)
- 코드 3-3 : [main.go](01-SimpleRPC/ordermgt/client/main.go)

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
