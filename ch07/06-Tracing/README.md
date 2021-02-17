# 오픈센서스 추적 (Chapter 7 : 서비스 수준 gRPC 실행)

## 예제 코드 리스트
- 코드 7-13 (예거 익스포터 초기화) : [tracer.go](productinfo/server/tracer/tracer.go)
- 코드 7-14 (서비스 계측 추가) : [main.go](productinfo/server/main.go)
- 코드 7-15 (클라아언트 계측 추가) : [main.go](productinfo/client/main.go)

## 1. 서비스 모니터링 활성화
다음과 같이 오픈센서스 기반 추적을 사용하기 위해 먼저 jaeger exporter를 초기화 합니다. [tracer.go](productinfo/server/tracer/tracer.go) (코드 7-13)
```go
package tracer

import (
	"log"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

func InitTracing() {
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: collectorEndpointURI,
		AgentEndpoint:     agentEndpointURI,
		ServiceName:       "product_info",
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
}
```
※ 도서 예제에는 `initTracing()` 함수명으로 정의되어 있으나, 외부 패키지에서 호출 시에는 `InitTracing()`와 같이 대문자로 시작해야 합니다.

그리고, 서비스 구현 부분에 다음과 같이 서비스 계측을 추가합니다. (main.go)[productinfo/server/main.go] (코드 7-14)

```go
// GetProduct는 ecommerce.GetProduct를 구현한다.
func (s *server) GetProduct(ctx context.Context, in *wrapper.StringValue) (
	*pb.Product, error) {
	ctx, span := trace.StartSpan(ctx, "ecommerce.GetProduct")
	defer span.End()
	value, exists := s.productMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}
```


## 2. 쿨라이언트 모니터링 활성화
클라이언트에서도 다음과 같이 오픈센서스 기반 추적을 사용해 모니터링을 활성화할 수 있습니다. [main.go](productinfo/client/main.go) (코드 7-15)
```go
package main

import (
	"context"
	"log"
	//"time" // 불필요

	pb "productinfo/client/ecommerce"

	wrapper "github.com/golang/protobuf/ptypes/wrappers" // 추가
	"productinfo/client/tracer"
	"google.golang.org/grpc"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	// "contrib.go.opencensus.io/exporter/jaeger" // 별도 tracer 패키지에서 사용
)

const (
	address = "localhost:50051"
)

func main() {
	tracer.InitTracing()	// initTracing()이 아닌 InitTracing()

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	ctx, span := trace.StartSpan(context.Background(),
		"ecommerce.ProductInfoClient")

	name := "Apple iphone 11"
	description := "Apple iphone 11 is the latest smartphone, launched in September 2019"
	price := float32(700.0)
	r, err := c.AddProduct(ctx, &pb.Product{Name: name,
		Description: description, Price: price})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	// Wrappers 사용
	//product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	product, err := c.GetProduct(ctx, &wrapper.StringValue{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: ", product.String())
	span.End()
}
```

## 3. 모니터링 정보 확인
서비스에서 jaeger로 정보를 보내기 때문에 우선 다음과 같이 jaeger를 docker 기반으로 실행합니다.

```shell
$ docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.21
```

그런 다음 다음과 같이 서비스 및 클라이언트를 실행합니다.
```shell
$ cd productinfo/server
$ go build -i -v -o bin/server main.go
$ bin/server &
$
$ cd ../../productinfo/client
$ go build -i -v -o bin/client main.go
$ bin/client
2021/02/17 16:38:42 Product ID: 288b146e-70f3-11eb-9bd6-acde48001122 added successfully
2021/02/17 16:38:42 Product: id:"288b146e-70f3-11eb-9bd6-acde48001122" name:"Apple iphone 11" description:"Apple iphone 11 is the latest smartphone, launched in September 2019" price:700
```

이제 다음과 같이 jaeger를 접속하여 추적 정보를 확인하십니다.
- http://localhost:16686/
