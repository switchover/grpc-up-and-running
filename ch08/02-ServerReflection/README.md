# 서버 리플렉션 (Chapter 8 : gRPC 생태계)

## 예제 코드 리스트
- 코드 8-3 (서비스 리플렉션 활성화) : [main.go](productinfo/server/main.go)

## 1. 서버 리플렉션 활성화
다음과 같이 서비스에 리플렉션을 지정합니다. [main.go](productinfo/server/main.go) (코드 8-3)
```go
package main

import (
	// ...

	pb "productinfo/server/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
// ...
```

## 2. 서비스 실행
리플렉션 확인을 위해 서비스를 다음과 같이 빌드 및 실행합니다.

```shell
$ cd productinfo/server
$ go build -i -v -o bin/server main.go
$ bin/server
```

## 3. gRPC CLI 실행
리플렉션 확인을 위해 우선 gRPC CLI(Command Line Interface)를 설치합니다. (참조 : https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md)

참고로 Mac의 경우는 다음과 같인 Homebrew를 통해 gRPC CLI를 설치할 수 있습니다.
```shell
$ brew install grpc
```

※ 참고로 도서에서는 직접 빌드한 경우로 현 디렉토리에 `grpc_cli` 실행 파일이 존재한다고 가정하였습니다. (아래 예시들은 Homebrew를 통해 설치된 경우로 설명하였습니다.)

### 1. 서비스 확인
```shell
$ grpc_cli ls localhost:50051
ecommerce.ProductInfo
grpc.reflection.v1alpha.ServerReflection
```

### 2. 서비스 상세 조회
```shell
$ grpc_cli ls localhost:50051 ecommerce.ProductInfo -l
filename: product_info.proto
package: ecommerce;
service ProductInfo {
  rpc addProduct(ecommerce.Product) returns (google.protobuf.StringValue) {}
  rpc getProduct(google.protobuf.StringValue) returns (ecommerce.Product) {}
}
```

### 3. 메서드 상세 조회
```shell
$ grpc_cli ls localhost:50051 ecommerce.ProductInfo.addProduct -l
rpc addProduct(ecommerce.Product) returns (google.protobuf.StringValue) {}
```

### 4. 메시지 타입 조회
```shell
$ grpc_cli type localhost:50051 ecommerce.Product
message Product {
  string id = 1 [json_name = "id"];
  string name = 2 [json_name = "name"];
  string description = 3 [json_name = "description"];
  float price = 4 [json_name = "price"];
}
```

### 5. 메서드 호출
```shell
$ grpc_cli call localhost:50051 addProduct "name:'Apple', description: 'iphone 11', price: 699"
connecting to localhost:50051
value: "2fc128f6-710a-11eb-bfbf-acde48001122"
Rpc succeeded with OK status
```

### 5. 메서드 호출
```shell
$ grpc_cli call localhost:50051 addProduct "name:'Apple', description: 'iphone 11', price: 699"
connecting to localhost:50051
value: "2fc128f6-710a-11eb-bfbf-acde48001122"
Rpc succeeded with OK status

$ grpc_cli call localhost:50051 getProduct "value:'2fc128f6-710a-11eb-bfbf-acde48001122'"
connecting to localhost:50051
id: "2fc128f6-710a-11eb-bfbf-acde48001122"
name: "Apple"
description: "iphone 11"
price: 699
Rpc succeeded with OK status
```