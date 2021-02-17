# Chapter 1 : gRPC 소개

## 예제 코드 리스트
- 코드 1-1 (서비스 정의 파일) : [ProductInfo.proto](ProductInfo.proto)
- 코드 1-2 (Go 서비스 구현): [productinfo/go/main.go](productinfo/go/main.go)
- 코드 1-3 (Go 서버 구현): [productinfo/go/main.go](productinfo/go/main.go)
- 코드 1-4 (Java 클라이언트 구현): [productinfo/java/src/main/java/Main.java](productinfo/java/src/main/java/Main.java)

## 정오
### 본문 부분
| 정오 | 수정 전 | 수정 후 |
| ------------- | ----- | ----- |
| [p.45 : 4행] | 기존 클라이언트-서버 통신에 근본적으로 다른 접근 방법의 고정된 계약 방식을 사용한다. | 기존의 클라이언트-서버 통신에 근본적으로 다른 접근 방식을 제공한다. 반면 gRPC는 클라이언트와 서버간 원격 메서드에 고정된 계약 방식을 사용한다 |

---
# `ProductInfo` 서비스 구현 (Go)

## 1. protobuf 정의 파일 생성
[ProductInfo.proto](ProductInfo.proto) (코드 1-1)

## 2. Go 서비스용 모듈 생성
우선 Go 모듈을 위한 디렉토리 생성 후, `go mod` 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
$ mkdir -p productinfo/go
$ cd productinfo/go
$ go mod init github.com/grpc-up-and-running/samples/ch02/productinfo/go
```
- 1장의 예제이지만, 책에 나와 있는 코드와 맞추기 위해서 `ch02` 디렉토리를 사용

## 3. protobuf 파일 복사
별도로 정의된 `ProductInfo.proto` 파일을 `proto` 디렉토리 생성 후 복사합니다.
```shell
$ mkdir proto
$ cp ../../ProductInfo.proto proto
```
- `ProductInfo.proto`는 임의 위치에서 복사함 (위 복사 경로는 현재 예제 디렉토리 구성의 예)

## 4. Go 언어 Skeleton 생성 
다음과 같이 이미 설치된 `protoc` 명령을 통해 skeleton 코드를 생성합니다.
```shell
$ protoc -I proto proto/ProductInfo.proto --go_out=plugins=grpc:proto 
```

## 5. Go 서비스 구현
다음과 같이 서비스를 구현합니다.
[main.go](productinfo/go/main.go) (코드 1-2)
```go
package main

import (
	"context"
	"errors"
	"log"
	"net"

	pb "github.com/grpc-up-and-running/samples/ch02/productinfo/go/proto"
	"google.golang.org/grpc"
)

type server struct {
	productMap map[string]*pb.Product
}

// Go 언어를 사용한 ProductInfo 구현

// 제품 등록을 위한 원격 메서드
func (s *server) AddProduct(ctx context.Context, in *pb.Product) (*pb.ProductID, error) {
	// 업무 로직
	in.Id = "Product1"

	log.Println("Name :", in.Name)
	log.Println("Desc. :", in.Description)

    if s.productMap == nil {
    	s.productMap = make(map[string]*pb.Product)
    }
	s.productMap[in.Id] = in

	return &pb.ProductID{Value: in.Id}, nil
}

// 제품 조회용 원격 메서드
func (s *server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.Product, error) {
	// 업무 로직

	product, exists := s.productMap[in.Value]
	if exists && product != nil {
		return product, nil
	}
	return nil, errors.New("product not found")
}
```

## Go main 함수 구현
Go의 main 함수를 다음과 같이 추가합니다.
[main.go](productinfo/go/main.go) (코드 1-3)
```go
const port = ":8080"

func main() {
	lis, _ := net.Listen("tcp", port)
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

## 6. Go 서버 빌드 및 실행
다음과 같이 서버를 빌드하고 실행합니다.
```shell
$ go build -i -v -o bin/server
$ bin/server
```

---
# `ProductInfo` 클라이언 구현 (Java)
## 1. Java 클라이언트 프로젝트 디렉토리 구조 생성
다음과 같이 gradle 프로젝트를 생성합니다.
```shell
$ mkdir productinfo/java
$ cd productinfo/java
$ mkdir src/main/java
$ mkdir src/main/proto
```

## 2. Gradle 빌드 파일 생성
다음과 같은 [build.gradle](productinfo/java/build.gradle) 파일을 생성합니다.

## 3. protobuf 파일 복사
별도로 정의된 `ProductInfo.proto` 파일을 gradle 프로젝트의 `src/main/proto` 디렉토리로 복사합니다.
```shell
$ cp ../../ProductInfo.proto src/main/proto
```
- `ProductInfo.proto`는 임의의 위치에서 복사함 (위 예는 현재 예제 디렉토리 구성의 경우임)

## 4. Java 언어 Stub 생성
다음과 같이 이미 설치된 `gradle` 명령을 통해 stub 코드를 생성합니다.
```shell
$ gradle build
```
- stub java 코드는 `build/generated/source/proto/main/grpc/` 및 `build/generated/source/proto/main/java/` 디렉토리 하위에 밑에 생성됨

## 5. Java 클라이언트 구현
다음과 같이 클라이언트를 구현합니다. 
[Main.java](productinfo/java/src/main/java/Main.java) (코드 1-4)
```java
import ecommerce.ProductInfoGrpc;
import ecommerce.ProductInfoOuterClass;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

public class Main {
    public static void main(String[] arg) {
        // 코드 1-4 부분

        // 원격 서버 주소를 사용해 채널(channel)을 생성
        ManagedChannel channel = ManagedChannelBuilder.forAddress("localhost", 8080)
                .usePlaintext(true)
                .build();
        // 채널을 사용해 블로킹(blocking) 방식 스텁 생성
        ProductInfoGrpc.ProductInfoBlockingStub stub;
        stub = ProductInfoGrpc.newBlockingStub(channel);

        // 블로킹 스텁을 통한 원격 메서드 호출
        ProductInfoOuterClass.ProductID productID = stub.addProduct(
                ProductInfoOuterClass.Product.newBuilder()
                        .setName("Apple iPhone 11")
                        .setDescription("Meet Apple iPhone 11." +
                                "All-new dual-camera system with " +
                                "Ultra Wide and Night mode.")
                        .build());
    }
}
```
- 책에서는 `stub.addProduct()` 메서드에 의해 리턴되는 `productID`의 타입이 `StringValue`로 되어 있으나
  정의된 proto 파일(`ProductInfo.proto`)는 일반 객체인 `ProductID`의 타입을 갖는다.
 
## 6. Java 클라이언트 빌드
다음과 같이 gradle을 통해 클라이언트를 다시 빌드합니다.
```shell
$ gradle build
```
- `build/libs/java.jar`로 빌드됨

## 7. Java 클라드이언트 실행
최종적으로 서버(`bin/server`)가 실행된 상태에서 다음과 같이 jar를 실행합니다.
```shell
$ java -jar build/libs/java.jar
```
- 호출된 결과는 서버로 console에서 다음과 같이 확인할 수 있습니다.
```shell
2021/02/03 00:09:35 Name : Apple iPhone 11
2021/02/03 00:09:35 Desc. : Meet Apple iPhone 11.All-new dual-camera system with Ultra Wide and Night mode.
```