# 프로메테우스 (Chapter 7 : 서비스 수준 gRPC 실행)

## 예제 코드 리스트
- 코드 7-10, 7-11 (서비스 모니터링 활성화) : [main.go](productinfo/server/main.go)
- 코드 7-12 (클라이언트 모니터링 활성화) : [main.go](productinfo/client/main.go)

## 1. 서비스 모니터링 활성화
다음과 같이 프로메테우스를 사용해 모니터링을 활성화할 수 있습니다. [main.go](productinfo/server/main.go) (코드 7-10)
```go
package main

import (
	// ...
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	reg = prometheus.NewRegistry()

	grpcMetrics = grpc_prometheus.NewServerMetrics()

	customMetricCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "product_mgt_server_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
)

func init() {
	reg.MustRegister(grpcMetrics, customMetricCounter)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9092)}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	pb.RegisterProductInfoServer(grpcServer, &server{})
	grpcMetrics.InitializeMetrics(grpcServer)

	// 프로메테우스에 대한 http 서버를 시작한다.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// ...
```

아울러, 개별 서비스 구현 부분에 다음과 같이 메트릭을 추가합니다. [main.go](productinfo/server/main.go) (코드 7-11)
```go
// AddProduct는 ecommerce.AddProduct를 구현한다.
func (s *server) AddProduct(ctx context.Context,
	in *pb.Product) (*wrapper.StringValue, error) {
	customMetricCounter.WithLabelValues(in.Name).Inc()
	// ...
```

## 2. 쿨라이언트 모니터링 활성화
클라이언트에서도 다음과 같이 프로메테우스를 사용해 모니터링을 활성화할 수 있습니다. [main.go](productinfo/client/main.go) (코드 7-12)
```go
package main

import (
	// ...
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	address = "localhost:50051"
)

func main() {
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	reg.MustRegister(grpcMetrics)

	conn, err := grpc.Dial(address,
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 프로메테우스용 HTTP 서버 생성
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9094)}

	// 프로메테우스용 HTTP 서버 시작
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	c := pb.NewProductInfoClient(conn)
	// ....
}
```

## 3. 모니터링 정보 확인
다음과 같이 모니터링 정보를 확인할 수 있다. 다만, 클라이언트는 경우는 실행 후 바로 종료되기 때문에 확인이 어렵습니다.
- http://localhost:9092/metrics (서비스)
- http://localhost:9094/metrics (클라이언트)
