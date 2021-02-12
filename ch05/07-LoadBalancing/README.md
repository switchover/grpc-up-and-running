# Section 7 : 로드밸런싱 (Chapter 5 : gRPC: 고급 기능)

## 예제 코드 리스트
- 코드 5-15 (네임 리졸버 구현), 코드 5-16 (클라이언트 로드밸런싱 구현): [main.go](echo/client/main.go)

----
# 로드밸런싱(Load Balancing) (Go)

## 1. Echo 서비스 확인
로드밸런싱을 테스트하기 위해 gRPC의 예제 중 [echo 서비스](https://github.com/grpc/grpc-go/blob/v1.29.1/examples/features/proto/echo/echo.proto)를 간단하게 구현하여 활용합니다.

`echo` 서비스에 대한 구현 코드는 [main.go](echo/server/main.go)를 참조하고, 
다음과 같이 빌드 및 실행합니다.

```shell
cd echo/server
go build -i -v -o bin/server main.go
bin/server
```
참고로 `echo` 서버는 2개의 서버를 시작합니다.


## 2. Name Resolver 구현
백엔드 IP의 목록을 반환하는 네임 리졸버를 다음과 같이 구현합니다.
먼저 네임 리졸버를 사용한 에코(echo) 서비스를 위한 모듈을 다음과 같이 구성합니다.
```shell
mkdir -p echo/client
cd echo/client
go mod init echo/client
```

실제로 구현된 네일 리졸버는 다음과 같이 구현됩니다. [main.go](echo/client/main.go) (코드 5-15)
```go
const (
	exampleScheme      = "example"
	exampleServiceName = "lb.example.grpc.io"
)
var addrs = []string{"localhost:50051", "localhost:50052"}

type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target,
	cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: addrs,
		},
	}
	r.start()
	return r, nil
}

func (*exampleResolverBuilder) Scheme() string { return exampleScheme }

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*exampleResolver) Close() {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
```

## 3. 클라이언트 로드밸런싱 구현
이미 구현된 네일 리졸버를 활용하여 다음과 같이 클라이언트 로드밸런싱을 구현합니다. [main.go](echo/client/main.go) (코드 5-16)
```go
pickfirstConn, err := grpc.Dial(
	fmt.Sprintf("%s:///%s",
		// exampleScheme = "example"
		// exampleServiceName = "lb.example.grpc.io"
		exampleScheme, exampleServiceName),
	// "pick_first"가 기본값이다.
	grpc.WithBalancerName("pick_first"),
	grpc.WithInsecure(),
)
if err != nil {
	log.Fatalf("did not connect: %v", err)
}
defer pickfirstConn.Close()

log.Println("==== Calling helloworld.Greeter/SayHello " +
	"with pick_first ====")
makeRPCs(pickfirstConn, 10)

// round_robin 정책으로 다른 ClientConn을 만든다.
roundrobinConn, err := grpc.Dial(
	fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
	// "example:///lb.example.grpc.io"
	grpc.WithBalancerName("round_robin"),
	grpc.WithInsecure(),
)
if err != nil {
	log.Fatalf("did not connect: %v", err)
}
defer roundrobinConn.Close()

log.Println("==== Calling helloworld.Greeter/SayHello " +
	"with round_robin ====")
makeRPCs(roundrobinConn, 10)
```
