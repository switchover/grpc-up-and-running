# 기타 인증 인증 (Chapter 6 : 보안 적용 gRPC)

## 1. 클라이언트 JWT 인증
JWT를 활용한 클라이언트 인증은 다음과 같이 구현합니다.
```go
jwtCreds, err := oauth.NewJWTAccessFromFile("token.json")
if err != nil {
	log.Fatalf("Failed to create JWT credentials: %v", err)
}

creds, err := credentials.NewClientTLSFromFile("server.crt", "localhost")
if err != nil {
	log.Fatalf("failed to load credentials: %v", err)
}
opts := []grpc.DialOption{
	grpc.WithPerRPCCredentials(jwtCreds),
	// 전송 자격증명(transport credentials)
	grpc.WithTransportCredentials(creds),
}

// 서버에 대한 연결을 설정한다.
conn, err := grpc.Dial(address, opts...)
if err != nil {
	log.Fatalf("did not connect: %v", err)
}
// .... // 스텁 생성과 RPC 메서드 호출은 생략한다.
```

## 2. 구글 클라우드 인증
구글 클라우드에 대한 클라인언트 인증은 다음과 같이 구현합니다.
```go
perRPC, err := oauth.NewServiceAccountFromFile("service-account.json", scope)
if err != nil {
	log.Fatalf("Failed to create JWT credentials: %v", err)
}

pool, _ := x509.SystemCertPool()
creds := credentials.NewClientTLSFromCert(pool, "")

opts := []grpc.DialOption{
	grpc.WithPerRPCCredentials(perRPC),
	grpc.WithTransportCredentials(creds),
}

conn, err := grpc.Dial(address, opts...)
if err != nil {
	log.Fatalf("did not connect: %v", err)
}
// .... // 스텁 생성과 RPC 메서드 호출은 생략한다.
```
