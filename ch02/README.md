# Chapter 2 : gRPC 시작

## 예제 코드 리스트
- 코드 2-5 : [product_info.proto](product_info.proto)
- 코드 2-6 : [productinfo_service.go](productinfo/service/productinfo_service.go)
- 코드 2-7 : [main.go](productinfo/service/main.go)
- 코드 2-8 : [build.gradle](product-info-service/build.gradle)
- 코드 2-9 : [ProductInfoImpl.java](product-info-service/src/main/java/ecommerce/ProductInfoImpl.java)
- 코드 2-10 : [ProductInfoServer.java](product-info-service/src/main/java/ecommerce/ProductInfoServer.java)
- 코드 2-11 : [product_client.go](productinfo/client/productinfo_client.go)
- 코드 2-12 : [ProductInfoClient.java](product-info-client/src/main/java/ecommerce/ProductInfoClient.java)

## 정오
### 본문 부분
| 정오 | 수정 전 | 수정 후 |
| ------------- | ----- | ----- |
| 54 페이지 하단 코드 예시 | 코드 2-4. | 코드 2-5. |

### 코드 부분
- 54, 55 페이지 코드 2-5. message 정의 부분 : `price` 정의 추가
    ```
    message Product {
        string id = 1;
        string name = 2;
        string description = 3;
    }
    ```
    :arrow_right:
    ```
    message Product {
        string id = 1;
        string name = 2;
        string description = 3;
        float price = 4;    // 추가
    }
    ```
    
- 61 페이지 코드 2-6. import 부분 : 사용 import 추가 및 불필요 import 제외
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


---
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

## 6. Go main 함수 구현
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

## 7. Go 서버 빌드 및 실행
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
이제 다음과 같이 gradle을 통해 서비스를 다시 빌드합니다.
```shell
gradle build
```
- `build/libs/product-info-service.jar`로 빌드됨

## 8. Java 서버 실행
다음과 같이 jar를 실행합니다.
```shell
java -jar build/libs/product-info-service.jar
```

---
# `ProductInfo` 클라이언트 구현 (Go)

## 1. Go 클라이언트용 모듈 생성
Go 모듈을 위한 디렉토리 생성 후, `go mod` 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
mkdir -p productinfo/client
cd productinfo/client
go mod init productinfo/client
```

## 2. protobuf 파일 복사
별도로 정의된 `product_info.proto` 파일을 `ecommerce` 디렉토리 생성 후 이 디렉토리로 복사합니다.
```shell
mkdir ecommerce
cp ../../product_info.proto ecommerce
```
- `product_info.proto`는 임의의 위치에서 복사함 (위 예는 현재 예제 디렉토리 구성의 경우임)

## 3. Go 언어 Stub 생성 
다음과 같이 이미 설치된 `protoc` 명령을 통해 skeleton 코드를 생성합니다.
```shell
protoc -I ecommerce ecommerce/product_info.proto --go_out=plugins=grpc:ecommerce 
```

## 4. Go 클라이언트 구현
다음과 같이 Go 클라이언트를 구현합니다.
[productinfo_client.go](productinfo/client/productinfo_client.go) (코드 2-11)
```go
package main

import (
	"context"
	"log"
	"time"

	pb "productinfo/client/ecommerce"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Apple iPhone 11"
	description := `Meet Apple iPhone 11. All-new dual-camera system with
        Ultra Wide and Night mode.`
    price := float32(1000.0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx,
		&pb.Product{Name: name, Description: description, Price: proce})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: ", product.String())
}
```

## 5. Go 클라이언트 빌드
다음과 같이 클라이언트를 빌드합니다.
```shell
go build -i -v -o bin/client
```

## 6. Go 클라이언트 실행
클라이언트를 실행하기 전에 Go 서버 또는 Java 서버가 실행되어 있는지 확인합니다.  
그런 다음 다음과 같이 Go 클라이언트를 실행합니다.
```shell
bin/client
```

---
# `ProductInfo` 클라이언트 구현 (Java)

## 1. Java 프로젝트 디렉토리 구조 생성
다음과 같이 gradle 프로젝트를 생성합니다.
```shell
mkdir product-info-client
cd product-info-client
mkdir -p src/main/java
mkdir -p src/main/proto
```

## 2. Gradle 빌드 파일 생성
다음과 같은 [build.gradle](product-info-client/build.gradle) 파일을 생성합니다.
- `product-info-service/build.gradle`에서 `Main-Class`만 변경 (`ecommerce.ProductInfoServer` -> `ecommerce.ProductInfoClient`)

## 3. protobuf 파일 복사
별도로 정의된 `ProductInfo.proto` 파일을 gradle 프로젝트의 `src/main/proto` 디렉토리로 복사합니다.
```shell
cp ../product_info.proto src/main/proto
```
- `ProductInfo.proto`는 임의 위치에서 복사함 (위 복사 경로는 현재 예제 디렉토리 구성의 예)

## 4. Java 언어 Stub 생성
다음과 같이 이미 설치된 `gradle` 명령을 통해 skeleton 코드를 생성합니다.
```shell
gradle build
```
- skeleton java 코드는 `build/generated/source/proto/main/grpc/` 및 `build/generated/source/proto/main/java/` 디렉토리 하위에 밑에 생성됨

## 5. Java 클라이언트 구현
이제 다음과 같이 클라이언트를 구현합니다. 
[ProductInfoClient.java](product-info-client/src/main/java/ecommerce/ProductInfoClient.java) (코드 2-12)
```java
package ecommerce;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import java.util.logging.Logger;

/**
* productInfo 서비스에 대한 gRPC 클라이언트 샘플
*/
public class ProductInfoClient {
    public static void main(String[] args) throws InterruptedException {
    ManagedChannel channel = ManagedChannelBuilder
        .forAddress("localhost", 50051)
        .usePlaintext()
        .build();

    ProductInfoGrpc.ProductInfoBlockingStub stub =
        ProductInfoGrpc.newBlockingStub(channel);

    ProductInfoOuterClass.ProductID productID = stub.addProduct(
        ProductInfoOuterClass.Product.newBuilder()
        .setName("Apple iPhone 11")
        .setDescription("Meet Apple iPhone 11. " +
            "All-new dual-camera system with " +
            "Ultra Wide and Night mode.")
        .setPrice(1000.0f)
        .build());
    System.out.println(productID.getValue());

    ProductInfoOuterClass.Product product = stub.getProduct(productID);
    System.out.println(product.toString());
    channel.shutdown();
    }
}
```
 
## 6. Java 클라이언트 빌드
이제 다음과 같이 gradle을 통해 클라이언트를 다시 빌드합니다.
```shell
gradle build
```
- `build/libs/product-info-client.jar`로 빌드됨

## 8. Java 클라이언트 실행
다음과 같이 jar를 실행합니다.
```shell
java -jar build/libs/product-info-client.jar
```