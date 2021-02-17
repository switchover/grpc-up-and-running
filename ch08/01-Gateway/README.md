# gRPC 게이트웨이 (Chapter 8 : gRPC 생태계)

## 예제 코드 리스트
- 코드 8-1 (서비스 정의) : [product_info.proto](proto/product_info.proto)
- 코드 8-2 (리버스 프록시 구현) : (main.go)[client/main.go]

## 1. 서비스 정의
다음과 같이 HTTP 리소스에 매핑 정보를 갖는 서비스를 정의합니다. [product_info.proto](proto/product_info.proto) (코드 8-1)
```
syntax = "proto3";

import "google/protobuf/wrappers.proto";
import "google/api/annotations.proto";

package ecommerce;

service ProductInfo {
    rpc addProduct(Product) returns (google.protobuf.StringValue) {
        option (google.api.http) = {
            post: "/v1/product"
            body: "*"
        };
    }
    rpc getProduct(google.protobuf.StringValue) returns (Product) {
         option (google.api.http) = {
             get:"/v1/product/{value}"
         };
    }
}

message Product {
    string id = 1;
    string name = 2;
    string description = 3;
    float price = 4;
}
```

## 2. Skeleton/Stub 소스 생성

이제 다음과 같이 Go Skeleton 및 Stub 코드를 생성합니다. (생성 전에 필요한 패키지를 다운로드힙니다.)

```shell
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2
$ go get -u github.com/golang/protobuf/protoc-gen-go

$ cd proto
$ protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--go_out=plugins=grpc:. \
product_info.proto
```
- (버전 변경사항)
    - `protoc-gen-swagger` 패키지가 `protoc-gen-openapiv2`로 변경됨에 따라 `go get` 명령이 일부 변경

아울러 다음과 같이 리버스 프록시용 코드도 생성합니다.
```shell
$ protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--grpc-gateway_out=logtostderr=true:. \
product_info.proto
```
- `product_info.pb.gw.go` 파일이 생성됨


## 3. 리버스 프록시 구현
생성된 코드를 활용하여 다음과 같이 reverse proxy를 구현합니다. (main.go)[client/main.go] (코드 8-2)

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

    gw "github.com/grpc-up-and-running/samples/ch08/grpc-gateway/go/gw"
)

var (
	grpcServerEndpoint = "localhost:50051"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterProductInfoHandlerFromEndpoint(ctx, mux,
		grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("Fail to register gRPC gateway service endpoint: %v", err)
	}

	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Could not setup HTTP endpoint: %v", err)
	}
}
```

## 4. 게이트웨이를 통한 호출
우선, 리버스 프록시를 실행하기 전에 서버를 먼저 빌드 및 실행합니다.

```shell
$ cd server
$ go build -i -v -o bin/server main.go
$ bin/server
```

다음으로 다음과 같이 리버스 프록시를 실행합니다.
```shell
$ cd client
$ go build -i -v -o bin/client main.go
$ bin/client
```

이제 curl을 통해 일반적인 HTTP를 호출하면 됩니다.
```shell
$ curl -X POST http://localhost:8081/v1/product \
-d '{"name": "Apple", "description": "iphone7", "price": 699}'
"9526ef66-7103-11eb-b484-acde48001122"
```

마지막 라인이 생성된 Product ID로 이 ID를 다음과 같이 사용해 제품 정보를 조회할 수 있습니다.
```shell
$ curl http://localhost:8081/v1/product/9526ef66-7103-11eb-b484-acde48001122
{"id":"9526ef66-7103-11eb-b484-acde48001122","name":"Apple","description":"iphone7","price":699}
```


## 5. Swagger 정의 파일 생성
gRPC Gateway는 다음과 같이 리버스 프록시 서비스의 스웨거 swagger 정의 파일도 생성할 수 있습니다.
```shell
$ protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/\
third_party/googleapis \
--openapiv2_out=logtostderr=true:. \
product_info.proto
```
- (버전 변경사항)
    - `--swagger_out=logtostderr=true:.` 부분이 `--openapiv2_out=logtostderr=true:.`로 변경 필요
- `product_info.swagger.json` 파일이 생성됨
