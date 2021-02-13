# mTLS 연결 (Chapter 6 : 보안 적용 gRPC)

## 예제 코드 리스트
- 코드 6-3 (서버 mTLS 보안 연결) : [main.go](productinfo/server/main.go)
- 코드 6-4 (클라이언트 mTLS 보안 연결) : [main.go](productinfo/client/main.go)

## 1. 추가 인증서 생성
mTLS 연결은 클라이언트에 대한 인증서 검증을 포함하며, 이를 위해서는 클라이언트 인증서도 필요합니다.

이에 대해서는 별도의 [인증서 생성](certs/Certificates.md) 문서를 참조합니다.

## 2. 서버 mTLS 보안 연결
서버에서는 다음과 같이 비밀키 및 인증서를 사용해 서버를 기동합니다. [main.go](productinfo/server/main.go) (코드 6-3)
```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	pb "productinfo/server/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"	
)

var (
	port    = ":50051"
	crtFile = "server.crt"
	keyFile = "server.key"
	caFile  = "ca.crt"
)

func main() {
	certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append ca certificate")
	}

	opts := []grpc.ServerOption{
		// 모든 요청 연결에 대해 TLS를 활성화한다.
		grpc.Creds(
			credentials.NewTLS(&tls.Config{
				ClientAuth:   tls.RequireAndVerifyClientCert,
				Certificates: []tls.Certificate{certificate},
				ClientCAs:    certPool,
			},
			)),
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
```

## 3. 클라이언트 mTLS 보안 연결
클라이언트에서는 다음과 같이 인증서를 사용해 연결을 시작합니다. [main.go](productinfo/client/main.go) (코드 6-4)
```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	pb "productinfo/server/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address  = "localhost:50051"
	hostname = "localhost"
	crtFile  = "client.crt"
	keyFile  = "client.key"
	caFile   = "ca.crt"
)

func main() {
	certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append ca certs")
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:   hostname, // 참고: 이 설정은 반드시 필요함!
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	// .... // RPC 메서드 호출 생략
```
