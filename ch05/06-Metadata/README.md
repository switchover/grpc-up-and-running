# Section 6 : 메타데이터 (Chapter 5 : gRPC: 고급 기능)

## 예제 코드 리스트
- 코드 5-11 (클라이언트 메타데이터 전송), 코드 5-12 (클라이언트 메타데이터 읽기) : [main.go](some-service/client/main.go)
- 코드 5-13 (서비스 메타데이터 읽기), 코드 5-14 (서비스 메타데이터 전송) [main.go](some-service/server/main.go)

----
# 메타데이터(Error Handling) (Go)

## 1. Some 서비스 정의
메타데이터 부분은 기존 주문 관리(Order Management) 서비스가 아닌 별도 서비스를 활용한 예제를 사용합니다.
따라서 별도 Some 서비스를 다음과 같이 정의합니다. (some.proto)[some.proto]

```
syntax = "proto3";

package some;

service Service {
    rpc SomeRPC(SomeRequest) returns (SomeResponse);
    rpc SomeStreamingRPC(stream SomeRequest) returns (SomeResponse);
}

message SomeRequest {
    string data = 1;
}

message SomeResponse {
    string data = 1;
}
```

## 2. Some 서비스용 모듈 생성 및 Go Skeleton 코드 생성
Go 모듈을 위한 디렉토리 생성 후, go mod 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
mkdir -p some-service/server
cd some-service/server
go mod init some-service/server

mkdir some
cp ../../some.proto some

protoc -I some some/some.proto --go_out=plugins=grpc:some 
```

## 3. Some 서비스 구현
정의된 Some 서비스에 대하여 다음과 같이 서버를 구현합니다. (main.go)[origin-some-service/server/main.go]

## 4. Some 클라이언트 모듈 생성 및 Go Stub 코드 생성
Go 모듈을 위한 디렉토리 생성 후, go mod 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
mkdir -p some-service/client
cd some-service/client
go mod init some-service/client

mkdir some
cp ../../some.proto some

protoc -I some some/some.proto --go_out=plugins=grpc:some 
```

## 5. Some 클라이언트 구현
정의된 Some 클라이언트에 대하여 다음과 같이 구현합니다. (main.go)[origin-some-service/client/main.go]

## 6. 클라이언트 메타데이터 전송 구현
다음과 같이 Go 클라이언트에서 메타데이터 전송을 적용합니다.
[main.go](some-service/client/main.go) (코드 5-11)

```go
md := metadata.Pairs(
	"timestamp", time.Now().Format(time.StampNano),
	"kn", "vn",
)
mdCtx := metadata.NewOutgoingContext(context.Background(), md)

ctxA := metadata.AppendToOutgoingContext(mdCtx,
	"k1", "v1", "k1", "v2", "k2", "v3")

// 단일 RPC 만들기
response, err := client.SomeRPC(ctxA, someRequest)
// ...

stream, err := client.SomeStreamingRPC(ctxA)
// ...
```

## 7. 클라이언트 메타데이터 읽기 구현
클라이언트에서도 다음과 같이 메타데이터를 읽을 수 있습니다.
[main.go](some-service/client/main.go) (코드 5-12)

```go
var header, trailer metadata.MD

// ***** 단일 RPC *****

r, err := client.SomeRPC(
	ctx,
	someRequest,
	grpc.Header(&header),
	grpc.Trailer(&trailer),
)

// 여기서 헤더와 트레일러 맵을 처리한다.

// ***** 스트리밍 RPC *****

stream, err := client.SomeStreamingRPC(ctx)

// 헤더를 조회
header, err := stream.Header()

// 트레일러 조회
trailer := stream.Trailer()

// 여기서 헤더와 트레일러 맵을 처리한다.
```

## 8. 서비스 메타데이터 읽기 구현
서버에서도 다음과 같이 메타데이터를 읽을 수 있습니다.
[main.go](some-service/server/main.go) (코드 5-13)

```go
func (s *server) SomeRPC(ctx context.Context,
	in *pb.SomeRequest) (*pb.SomeResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	// 메타데이터 활용

// ...

func (s *server) SomeStreamingRPC(
	stream pb.Service_SomeStreamingRPCServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	// 메타데이터 활용

// ...
```

## 9. 서비스 메타데이터 전송 구현
서버에서도 다음과 같이 메타데이터를 전송합니다.
[main.go](some-service/server/main.go) (코드 5-14)

```go
func (s *server) SomeRPC(ctx context.Context,
	in *pb.someRequest) (*pb.someResponse, error) {
	// 헤더 생성과 전송
	header := metadata.Pairs("header-key", "val")
	grpc.SendHeader(ctx, header)
	// 트레일러 생성과 지정
	trailer := metadata.Pairs("trailer-key", "val")
	grpc.SetTrailer(ctx, trailer)
}

func (s *server) SomeStreamingRPC(stream pb.Service_SomeStreamingRPCServer) error {
	// 헤더 생성과 전송
	header := metadata.Pairs("header-key", "val")
	stream.SendHeader(header)
	// 트레일러 생성과 지정
	trailer := metadata.Pairs("trailer-key", "val")
	stream.SetTrailer(trailer)
}
```
