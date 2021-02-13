# OAuth 2.0 인증 (Chapter 6 : 보안 적용 gRPC)

## 예제 코드 리스트
- 코드 6-8 (클라이언트 OAuth 인증 처리) : [main.go](productinfo/client/main.go)
- 코드 6-9 (서버 OAuth 인증처리) : [main.go](productinfo/server/main.go)

## 1. 클라이언트 OAuth 인증 처리
클라이언트에서는 다음과 같이 베이직 인증을 처리합니다. [main.go](productinfo/client/main.go) (코드 6-6)
```go
package main

import (
	"context"
	"log"
	"time"

	pb "productinfo/client/ecommerce"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	crtFile  = "server.crt"
)

func main() {
	auth := oauth.NewOauthAccess(fetchToken())

	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
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

	// .... // RPC 메소드 호출을 건너뛴다.

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}
```

## 2. 서버 OAuth 인증 처리
서버에서는 다음과 같이 베이직 인증을 처리합니다. [main.go](productinfo/server/main.go) (코드 6-7)
```go
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"strings"
	
	pb "productinfo/server/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// server는 ecommerce/product_info를 구현하는 데 사용한다.
type server struct {
	productMap map[string]*pb.Product
}

var (
	port               = ":50051"
	crtFile            = "server.crt"
	keyFile            = "server.key"
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func main() {
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(ensureValidToken),
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
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "some-secret-token"
}

func ensureValidToken(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}
```
