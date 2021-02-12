# Section 5 : 멀티플렉싱 (Chapter 5 : gRPC: 고급 기능)

## 예제 코드 리스트
- 코드 5-9 (서버 멀티플렉싱) : [main.go](order-service/server/main.go)
- 코드 5-10 (클라이언트 멀티플렉싱) : [main.go](order-service/client/main.go)

----
# 멀티플렉싱(Multiplexing) (Go)

## 1. 추가 Hello Service 
기준 주문 관리 서비스((Order Management)[../order_management.proto])외에 추가로 Hello World 서비스를 다음과 같이 정의합니다. (helloworld.proto)[helloworld.proto]

```
syntax = "proto3";

package helloworld;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

다음으로 기존 주문 관리 서비스에 Hello World 서비스에 대한 Go Skeleton 및 Go Stub 코드를 다음과 같이 생성합니다.
```shell
mkdir order-service/server/helloworld
cp helloworld.proto order-service/server/helloworld
protoc -I order-service/server/helloworld order-service/server/helloworld/helloworld.proto --go_out=plugins=grpc:order-service/server/helloworld

mkdir order-service/client/helloworld
cp helloworld.proto order-service/client/helloworld
protoc -I order-service/client/helloworld order-service/client/helloworld/helloworld.proto --go_out=plugins=grpc:order-service/client/helloworld
```

## 2. 서버 멀티플렉싱(Server multiplexing) 구현
다음과 같이 Go 서비스에 멀티플렉싱을 적용합니다.
[main.go](order-service/server/main.go) (코드 5-9)

```go
func main() {
	initSampleData()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// gRPC orderMgtServer에 주문 관리 서비스 등록
	ordermgt_pb.RegisterOrderManagementServer(grpcServer, &orderMgtServer{})

	// gRPC orderMgtServer에 Greeter 서비스 등록
	hello_pb.RegisterGreeterServer(grpcServer, &helloServer{})

	// ...
```

## 3. 클라이언트 멀티플렉싱(Client multiplexing) 구현
다음과 같이 Go 클라이언트에 멀티플렉싱을 적용합니다.
[main.go](order-service/client/main.go) (코드 5-10)

```go
// 서버에 대한 연결을 설정한다.
conn, err := grpc.Dial(address, grpc.WithInsecure())
// ...

orderManagementClient := pb.NewOrderManagementClient(conn)
// ...

// Add Order
// ...
res, addErr := orderManagementClient.AddOrder(ctx, &order1)
// ...

helloClient := hwpb.NewGreeterClient(conn)
// ...

// Say hello RPC
helloResponse, err := helloClient.SayHello(hwcCtx,
	&hwpb.HelloRequest{Name: "gRPC Up and Running!"})
// ...
```
