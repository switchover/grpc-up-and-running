# Chapter 2 : gRPC 시작

## 예제 코드 리스트
- 코드 2-5 : [product_info.proto](product_info.proto)
- 코드 2-6 : [productinfo_service.go](productinfo/service/productinfo_service.go)
- 코드 2-7 : [main.go](productinfo/service/main.go)
- 코드 2-8 : [build.gradle]
- 코드 2-9 : [ProductInfoImpl.java](product-info-service/src/main/java/ecommerce/ProductInfoImpl.java)
- 코드 2-10 : [ProductInfoServer.java](product-info-service/src/main/java/ecommerce/ProductInfoServer.java)

## 정오
### 본문 부분
| 정오 | 수정 전 | 수정 후 |
| ------------- | ----- | ----- |
| 54 페이지 하단 코드 예시 | 코드 2-4. | 코드 2-5. |

### 코드 부분
- 61 페이지 코드 2-6. import 부분 
    ```go
    import (
        "context"
        "errors"
        "log"
        "github.com/gofrs/uuid"
        pb "productinfo/service/ecommerce"
    )
    ```
    :arrow_right:
    ```go
    import (
        "context"

        pb "productinfo/service/ecommerce"

        "github.com/gofrs/uuid"
        "google.golang.org/grpc/codes"
        "google.golang.org/grpc/status"
    )
    ```


# `ProductInfo` 서비스 구현 (Go)

## 1. protobuf 정의 파일 생성
[product_info.proto](product_info.proto) (코드 2-5)

## 2. Go 서비스용 모듈 생성
Go 모듈을 위한 디렉토리 생성 후, `go mod` 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
mkdir -p productinfo/service
cd productinfo/service
go mod init productinfo/service
```

## 3. protobuf 파일 복사
별도로 정의된 `product_info.proto` 파일을 `ecommerce` 디렉토리 생성 후 이 디렉토리로 복사합니다.
```shell
mkdir ecommerce
cp ../../product_info.proto ecommerce
```
- `product_info.proto`는 임의의 위치에서 복사함 (위 예는 현재 예제 디렉토리 구성의 경우임)

## 4. Go 언어 Skeleton 생성 
다음과 같이 이미 설치된 `protoc` 명령을 통해 skeleton 코드를 생성합니다.
```shell
protoc -I ecommerce ecommerce/product_info.proto --go_out=plugins=grpc:ecommerce 
```

## 5. Go 서비스 구현
다음과 같이 Go 서비스를 구현합니다.
[productinfo_service.go](productinfo/service/productinfo_service.go) (코드 2-6)
```go
// 코드 2-6 부분
package main

import (
	"context"

	pb "productinfo/service/ecommerce"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server는 ecommerce/product_info를 구현하는 데 사용된다.
type server struct {
	productMap map[string]*pb.Product
}

// AddProduct는 ecommerce.AddProduct를 구현한다.
func (s *server) AddProduct(ctx context.Context,
	in *pb.Product) (*pb.ProductID, error) {
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"Error while generating Product ID", err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[in.Id] = in
	return &pb.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

// GetProduct는 ecommerce.GetProduct를 구현한다.
func (s *server) GetProduct(ctx context.Context,
	in *pb.ProductID) (*pb.Product, error) {

	value, exists := s.productMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}
```
- (정오) 책에서의 `import` 처리에는 `google.golang.org/grpc/codes` 패키지 등에 대한 import가 누락되고, 사용되지 않는 `errors`, `log` 패키지가 포함되어 있다.

## Go main 함수 구현
이제 마지막으로 Go의 main 함수룰 다음과 같이 추가합니다.
[main.go](productinfo/service/main.go) (코드 2-7)
```go
package main

import (
	"log"
	"net"

	pb "productinfo/service/ecommerce"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

## 6. Go 서버 빌드 및 실행
다음과 같이 서버를 빌드하고 실행합니다.
```shell
go build -i -v -o bin/server
bin/server
```

---
# `ProductInfo` 서비스 구현 (Java)
## 1. Java 프로젝트 디렉토리 구조 생성
다음과 같이 gradle 프로젝트를 생성합니다.
```shell
mkdir product-info-service
cd product-info-service
mkdir -p src/main/java
mkdir -p src/main/proto
```

## 2. Gradle 빌드 파일 생성
다음과 같은 [build.gradle](product-info-service/build.gradle) 파일을 생성합니다. (코드 2-8)

## 3. protobuf 파일 복사
별도로 정의된 `ProductInfo.proto` 파일을 gradle 프로젝트의 `src/main/proto` 디렉토리로 복사합니다.
```shell
cp ../product_info.proto src/main/proto
```
- `ProductInfo.proto`는 임의 위치에서 복사함 (위 복사 경로는 현재 예제 디렉토리 구성의 예)

## 4. Java 언어 Skeleton 생성
다음과 같이 이미 설치된 `gradle` 명령을 통해 skeleton 코드를 생성합니다.
```shell
gradle build
```
- skeleton java 코드는 `build/generated/source/proto/main/grpc/` 및 `build/generated/source/proto/main/java/` 디렉토리 하위에 밑에 생성됨

## 5. Java 서비스 구현
이제 다음과 같이 서비스를 구현합니다. 
[ProductInfoImpl.java](product-info-service/src/main/java/ecommerce/ProductInfoImpl.java) (코드 2-9)
```java
package ecommerce;

import io.grpc.Status;
import io.grpc.StatusException;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

public class ProductInfoImpl extends ProductInfoGrpc.ProductInfoImplBase {
    private Map productMap = new HashMap<String, ProductInfoOuterClass.Product>();

    @Override
    public void addProduct(
        ProductInfoOuterClass.Product request,
        io.grpc.stub.StreamObserver
            <ProductInfoOuterClass.ProductID> responseObserver) {
        UUID uuid = UUID.randomUUID();
        String randomUUIDString = uuid.toString();
        request = request.toBuilder().setId(randomUUIDString).build();
        productMap.put(randomUUIDString, request);
        ProductInfoOuterClass.ProductID id =
            ProductInfoOuterClass.ProductID.newBuilder()
            .setValue(randomUUIDString).build();
            responseObserver.onNext(id);
        responseObserver.onCompleted();
    }

    @Override
    public void getProduct(
        ProductInfoOuterClass.ProductID request,
        io.grpc.stub.StreamObserver
            <ProductInfoOuterClass.Product> responseObserver) {
        String id = request.getValue();
        if (productMap.containsKey(id)) {
            responseObserver.onNext(
                (ProductInfoOuterClass.Product) productMap.get(id));
            responseObserver.onCompleted();
        } else {
            responseObserver.onError(new StatusException(Status.NOT_FOUND));
        }
    }
}
```
 
## 6. Java 서버 구현
Java 서버를 다음과 같이 구현합니다.
[ProductInfoServer.java](product-info-service/src/main/java/ecommerce/ProductInfoServer.java) (코드 2-10)
```java
package ecommerce;

import io.grpc.Server;
import io.grpc.ServerBuilder;

import java.io.IOException;

public class ProductInfoServer {
    public static void main(String[] args)
            throws IOException, InterruptedException {
        int port = 50051;
        Server server = ServerBuilder.forPort(port)
            .addService(new ProductInfoImpl())
            .build()
            .start();
        System.out.println("Server started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.err.println("Shutting down gRPC server since JVM is " +
                "shutting down");
            if (server != null) {
                server.shutdown();
            }
            System.err.println("Server shut down");
        }));
        server.awaitTermination();
    }
}
```

## 7. Java 서비스 빌드
이제 다음과 같이 gradle을 통해 서비를 다시 빌드합니다.
```shell
gradle build
```
- `build/libs/java.jar`로 빌드됨

## 8. Java 서버 실행
다음과 같이 jar를 실행합니다.
```shell
java -jar build/libs/product-info-service.jar
```

---
# `ProductInfo` 클라이언트 구현 (Go)

TBD

# `ProductInfo` 클라이언트 구현 (Java)

TBD
