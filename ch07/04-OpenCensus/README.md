# 오픈센서스 (Chapter 7 : 서비스 수준 gRPC 실행)

## 예제 코드 리스트
- 코드 7-8 (서비스 모니터링 활성화) : [main.go](productinfo/server/main.go)
- 코드 7-9 (클라이언트 모니터링 활성화) : [main.go](productinfo/client/main.go)

## 1. 서비스 모니터링 활성화
다음과 같이 오픈센서스를 사용해 모니터링을 활성화할 수 있습니다. [main.go](productinfo/server/main.go) (코드 7-8)
```go
package main

import (
  "errors"
  "log"
  "net"
  "net/http"

  pb "productinfo/server/ecommerce"
  "google.golang.org/grpc"
  "go.opencensus.io/plugin/ocgrpc"
  "go.opencensus.io/stats/view"
  "go.opencensus.io/zpages"
  "go.opencensus.io/examples/exporter"
)

var (
	port = ":50051"
)

// ecommerce/product_info를 구현하는 서버
type server struct {
	productMap map[string]*pb.Product
}

func main() {
	go func() {
		mux := http.NewServeMux()
		zpages.Handle(mux, "/debug")
		log.Fatal(http.ListenAndServe("127.0.0.1:8081", mux))
	}()

	view.RegisterExporter(&exporter.PrintExporter{})

	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	pb.RegisterProductInfoServer(grpcServer, &server{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```
참고로 도서의 예제는 서비스의 구현 부분은 제외되어 있으며, 실제 서비스 구현을 위해서는 `import`가 추가되어야 합니다. ([main.go](productinfo/server/main.go) 참조)

## 2. 쿨라이언트 모니터링 활성화
클라이언트에서도 다음과 같이 오픈센서스를 사용해 모니터링을 활성화할 수 있습니다. [main.go](productinfo/client/main.go) (코드 7-9)
```go
package main

import (
	"context"
	"log"
	"time"

  pb "productinfo/server/ecommerce"
  "google.golang.org/grpc"
  "go.opencensus.io/plugin/ocgrpc"
  "go.opencensus.io/stats/view"
  "go.opencensus.io/examples/exporter"
)

const (
	address = "localhost:50051"
)

func main() {
	view.RegisterExporter(&exporter.PrintExporter{})

	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(address,
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Can't connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewProductInfoClient(conn)
	// .... // RPC 메서드 호출 생략
```

## 3. 모니터링 정보 확인
서비스 실행 후, 클라이언트를 실행하면 서비스 console 상에 다음과 같은 출력 정보가 표시됩니다.
```shell
13:52:03 grpc.io/server/sent_messages_per_rpc         distribution: min=1.0 max=1.0 mean=1.0
  - grpc_server_method=ecommerce.ProductInfo/addProduct
13:52:03 grpc.io/server/sent_messages_per_rpc         distribution: min=1.0 max=1.0 mean=1.0
  - grpc_server_method=ecommerce.ProductInfo/getProduct
13:52:03 grpc.io/server/server_latency                distribution: min=0.9 max=0.9 mean=0.9
  - grpc_server_method=ecommerce.ProductInfo/addProduct
13:52:03 grpc.io/server/server_latency                distribution: min=0.1 max=0.1 mean=0.1
  - grpc_server_method=ecommerce.ProductInfo/getProduct
13:52:03 grpc.io/server/completed_rpcs                count:        value=1
  - grpc_server_method=ecommerce.ProductInfo/getProduct
  - grpc_server_status=OK
13:52:03 grpc.io/server/completed_rpcs                count:        value=1
  - grpc_server_method=ecommerce.ProductInfo/addProduct
  - grpc_server_status=OK
13:52:03 grpc.io/server/received_bytes_per_rpc        distribution: min=87.0 max=87.0 mean=87.0
  - grpc_server_method=ecommerce.ProductInfo/addProduct
13:52:03 grpc.io/server/received_bytes_per_rpc        distribution: min=38.0 max=38.0 mean=38.0
  - grpc_server_method=ecommerce.ProductInfo/getProduct
13:52:03 grpc.io/server/received_messages_per_rpc     distribution: min=1.0 max=1.0 mean=1.0
  - grpc_server_method=ecommerce.ProductInfo/getProduct
13:52:03 grpc.io/server/received_messages_per_rpc     distribution: min=1.0 max=1.0 mean=1.0
  - grpc_server_method=ecommerce.ProductInfo/addProduct
13:52:03 grpc.io/server/sent_bytes_per_rpc            distribution: min=38.0 max=38.0 mean=38.0
  - grpc_server_method=ecommerce.ProductInfo/addProduct
13:52:03 grpc.io/server/sent_bytes_per_rpc            distribution: min=125.0 max=125.0 mean=125.0
  - grpc_server_method=ecommerce.ProductInfo/getProduct
```

아울러, 다음과 같은 URL을 통해서 모니터링 정보를 확인할 수 있습니다.
- http://127.0.0.1:8081/debug/rpcz
- http://127.0.0.1:8081/debug/tracez

