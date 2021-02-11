# Section 1 : 인터셉터 (Chapter 5 : gRPC: 고급 기능)

## 예제 코드 리스트
- 코드 5-1 (서버 단일 인터셉터), 코드 5-2 (서버 스트리밍 인터셉터) : [main.go](order-service/server/main.go)
- 코드 5-3 (클라이언트 단일 인터셉터), 코드 5-4 (클라이언트 스트리밍 인터셉터) : [main.go](order-service/client/main.go)

----
# 인터셉터(Interceptors) (Go)

## 1. 서버 단일 인터셉터(Server-side unary interceptor) 구현
다음과 같이 Go 서비스에 서버측 단일 인터셉터를 추가합니다.
[main.go](order-service/server/main.go)

```go
// 서버 - 단일 인터셉터
func orderUnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// 전처리 로직
	// 인자(args)로 넘겨진 info를 통해 현재 RPC 호출에 대한 정보를 얻는다.
	log.Println("======= [Server Interceptor] ", info.FullMethod)

	// 단일 RPC의 정상 실행을 완료하고자 핸들러(handler)를 호출한다.
	m, err := handler(ctx, req)

	// 후처리 로직
	log.Printf(" Post Proc Message : %s", m)
	return m, err
}

func main() {
    // ...
    // 서버 측에서 인터셉터를 등록한다.
	s := grpc.NewServer(
		grpc.UnaryInterceptor(orderUnaryServerInterceptor))
```

## 2. 서버 스트리밍 인터셉터(Server-side streaming interceptor) 구현
다음과 같이 Go 서비스에 서버측 스트리밍 인터셉터를 추가합니다.
[main.go](order-service/server/main.go)

```go
// 서버 - 스트리밍 인터셉터
// wrappedStream이 내부의 grpc.ServerStream을 감싸고,
// RecvMsg와 SendMsg 메서드 호출을 가로챈다.
type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("====== [Server Stream Interceptor Wrapper] "+
		"Receive a message (Type: %T) at %s",
		m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("====== [Server Stream Interceptor Wrapper] "+
		"Send a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func orderServerStreamInterceptor(srv interface{},
	ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	log.Println("====== [Server Stream Interceptor] ",
		info.FullMethod)
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return err
}

func main() {
    // ...
    // 인터셉터 등록
    s := grpc.NewServer(
        grpc.StreamInterceptor(orderServerStreamInterceptor))
```


## 3. 클라이언트 단일 인터셉터(Client-side unary interceptor) 구현
다음과 같이 Go 클라이언트에 클라이언트측 단일 인터셉터를 추가합니다.
[main.go](order-service/client/main.go)

```go
func orderUnaryClientInterceptor(
	ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 전처리 단계
	log.Println("Method : " + method)

	// 원격 메서드 호출
	err := invoker(ctx, method, req, reply, cc, opts...)

	// 후처리 단계
	log.Println(reply)

	return err
}

func main() {
    // 서버로의 연결을 설정한다.
    conn, err := grpc.Dial(address, grpc.WithInsecure(),
        grpc.WithUnaryInterceptor(orderUnaryClientInterceptor))
    // ...
```


## 4. 클라이언트 스트리밍 인터셉터(Client-side streaming interceptor) 구현
다음과 같이 Go 클라이언트에 클라이언트측 스트리밍 인터셉터를 추가합니다.
[main.go](order-service/client/main.go)

```go
func clientStreamInterceptor(
	ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string,
	streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Println("======= [Client Interceptor] ", method)
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStream(s), nil
}

type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("====== [Client Stream Interceptor] "+
		"Receive a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("====== [Client Stream Interceptor] "+
		"Send a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}


func main() {
    // 서버와의 연결을 구성한다.
    conn, err := grpc.Dial(address, grpc.WithInsecure(),
        grpc.WithStreamInterceptor(clientStreamInterceptor))
    // ...
```