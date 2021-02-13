# 베이직 인증 (Chapter 6 : 보안 적용 gRPC)

## 예제 코드 리스트
- 코드 6-5 (PerRPCCredentials 인터페이스 구현), 코드 6-6 (클라이언트 베이직 인증 처리) : [main.go](productinfo/client/main.go)
- 코드 6-7 (서버 베이직 인증처리) : [main.go](productinfo/server/main.go)

## 1. 베이직 인증 인터페잇 구현
베이직 인증을 위해서는 `PerRPCCredentials` 인터페이스를 다음과 같이 구현해야 합니다. [main.go](productinfo/client/main.go) (코드 6-5)
```go
type basicAuth struct {
	username string
	password string
}

func (b basicAuth) GetRequestMetadata(ctx context.Context,
	in ...string) (map[string]string, error) {
	auth := b.username + ":" + b.password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

func (b basicAuth) RequireTransportSecurity() bool {
	return true
}
```

## 2. 클라이언트 베이직 인증 처리
클라이언트에서는 다음과 같이 베이직 인증을 처리합니다. [main.go](productinfo/client/main.go) (코드 6-6)
```go
package main

import (
	"log"
	pb "productinfo/server/ecommerce"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	crtFile  = "server.crt"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	auth := basicAuth{
		username: "admin",
		password: "admin",
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(auth),
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewProductInfoClient(conn)

	// .... // RPC 메서드 호출 생략
```

## 3. 서버 베이직 인증 처리
서버에서는 다음과 같이 베이직 인증을 처리합니다. [main.go](productinfo/server/main.go) (코드 6-7)
```go
package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	pb "productinfo/server/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"path/filepath"
	"strings"
)

var (
	port               = ":50051"
	crtFile            = "server.crt"
	keyFile            = "server.key"
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid credentials")
)

type server struct {
	productMap map[string]*pb.Product
}

func main() {
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	opts := []grpc.ServerOption{
		// 모든 들어오는 연결에 대해 TLS를 활성화한다.
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(ensureValidBasicCredentials),
	}

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

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Basic ")
	return token == base64.StdEncoding.EncodeToString([]byte("admin:admin"))
}

func ensureValidBasicCredentials(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// 유효한 토큰을 확인한 후 핸들러 실행을 계속한다.
	return handler(ctx, req)
}
```
