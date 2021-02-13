# 보안 설정 (Chapter 6 : 보안 적용 gRPC)

## 예제 코드 리스트
- 코드 6-1 (서버 단방향 보안 연결) : [main.go](productinfo/server/main.go)
- 코드 6-2 (클라이언트 단방향 보안 연결) : [main.go](productinfo/client/main.go)

## 1. 서버 단방향(One-Way) 보안 연결
서버에서는 다음과 같이 비밀키 및 인증서를 사용해 서버를 기동합니다. [main.go](productinfo/server/main.go) (코드 6-1)
```go
package main

import (
	"crypto/tls"
	"errors"
    pb "productinfo/server/ecommerce"
    "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

var (
	port    = ":50051"
	crtFile = "server.crt"
	keyFile = "server.key"
)

func main() {
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)), // 여러 줄인 경우 "," 추가 필요
	}
	// 또는 다음과 같이 사용 가능
	// opts := []grpc.ServerOption{grpc.Creds(credentials.NewServerTLSFromCert(&cert))}

	s := grpc.NewServer(opts...)

	pb.RegisterProductInfoServer(s, &server{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

## 2. 클라이언트 단방향(One-Way) 보안 연결
클라이언트에서는 다음과 같이 인증서를 사용해 연결을 시작합니다. [main.go](productinfo/client/main.go) (코드 6-2)
```go
package main

import (
    "log"
    pb "productinfo/server/ecommerce"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc"
)

const (
	address  = "localhost:50051"
	hostname = "localhost"
	crtFile  = "server.crt"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	// ... // RPC 메서드 호출 생략
```