# 전체 정오표

---
## 1장
### 본문 부분
| 정오 | 수정 전 | 수정 후 |
| ------------- | ----- | ----- |
| [p.45 : 4행] | 기존 클라이언트-서버 통신에 근본적으로 다른 접근 방법의 고정된 계약 방식을 사용한다. | 기존의 클라이언트-서버 통신에 근본적으로 다른 접근 방식을 제공한다. 반면 gRPC는 클라이언트와 서버간 원격 메서드에 고정된 계약 방식을 사용한다 |

---
## 2장
### 본문 부분
| 정오 | 수정 전 | 수정 후 |
| ------------- | ----- | ----- |
| [p.54 하단 코드 예시 캡션] | 코드 2-4. | 코드 2-5. |
| [p.66 옮긴이 메모] | 코드 2-4. | 코드 2-5. |

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
        float price = 4;
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
## 3장
### 코드 부분
- 98 페이지 코드 3-9. `client` 정의 부분 : 변수명 변경
    ```go
    c := pb.NewOrderManagementClient(conn)
    ```
    :arrow_right:
    ```go
    client := pb.NewOrderManagementClient(conn)
    ```
- 106 페이지 코드 3-12. `asyncClientBidirectionalRPC()` 함수 마지막 부분 : channel 종료 처리
    ```go
    <-c
    ```
    :arrow_right:
    ```go
    close(c)
    ```

---
## 4장
해당 없음

---
## 5장
### 코드 부분
- 173 페이지 코드 5-15. `Build` 메서드 파라미터 타입 부분 : `resolver.BuildOption` -> `resolver.BuildOptions`
    ```go
    func (*exampleResolverBuilder) Build(target resolver.Target,
        cc resolver.ClientConn,
        opts resolver.BuildOption) (resolver.Resolver, error) {
    ```
    :arrow_right:
    ```go
    func (*exampleResolverBuilder) Build(target resolver.Target,
        cc resolver.ClientConn,
        opts resolver.BuildOptions) (resolver.Resolver, error) {
    ```
- 174 페이지 코드 5-15. `ResolveNow` 메서드 파라미터 타입 부분 : `resolver.ResolveNowOption` -> `resolver.ResolveNowOptions`
    ```go
    func (*exampleResolver) ResolveNow(o resolver.ResolveNowOption) {}
    ```
    :arrow_right:
    ```go
    func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
    ```

---
## 6장
### 코드 부분
- 185 페이지 코드 6-1. `grpc.ServerOption` 구조체 리터럴 부분 : 뒤에 `,` 추가 필요 (여러 줄을 사용한 경우)
    ```go
    opts := []grpc.ServerOption{
        grpc.Creds(credentials.NewServerTLSFromCert(&cert))
    }
    ```
    :arrow_right:
    ```go
    opts := []grpc.ServerOption{
        grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
    }
    ```
    또는
    ```go
    opts := []grpc.ServerOption{grpc.Creds(credentials.NewServerTLSFromCert(&cert))}
    ```
- 187 페이지 코드 6-2. `var` 선언 중 `hostname` 지정 부분 : 뒤에 `"` 누락
    ```go
    address = "localhost:50051
    ```
    :arrow_right:
    ```go
    address = "localhost:50051"
    ```

---
## 7장
### 코드 부분
- 217 페이지 ghz 실행 : IP(`0.0.0.0:50051`) 부분  
    ```shell
    0.0.0.0:50051
    ```
    :arrow_right:
    ```shell
    localhost:50051
    ```
- 220 페이지 docker image build 실행 : 뒤 부분 경로 추가 필요  
    ```shell
    docker image build -t grpc-productinfo-server -f server/Dockerfile
    ```
    :arrow_right:
    ```shell
    docker image build -t grpc-productinfo-server -f server/Dockerfile .
    ```
- 244 페이지 코드 7-13 (선택사항) : `initTracing()` 함수명 변경(`InitTracing()`) => 별도 패키지로 정의되어 있어 다른 패키지에서 사용 시에는 export name으로 정의(대문자 시작) 필요
    ```go
    func initTracing() {
    ```
    :arrow_right:
    ```go
    func InitTracing() {
    ```
- 245 페이지 코드 7-15 import 부분 : tracer 패키지 분리에 따라 불필요한 `import` 정리 및 `wrapper` 추가 
    ```go
    import (
        "context"
        "log"
        "time"
        pb "productinfo/client/ecommerce"
        "productinfo/client/tracer"
        "google.golang.org/grpc"
        "go.opencensus.io/plugin/ocgrpc"
        "go.opencensus.io/trace"
        "contrib.go.opencensus.io/exporter/jaeger"
    )
    ```
    :arrow_right:
    ```go
    import (
        "context"
        "log"

        pb "productinfo/client/ecommerce"
        "productinfo/client/tracer"

        wrapper "github.com/golang/protobuf/ptypes/wrappers"
        "go.opencensus.io/trace"
        "google.golang.org/grpc"
    )
    ``` 
- 246 페이지 코드 7-15 `tracer.initTrace()` 호출 부분 : 함수명 변경(`InitTracing()`) => 별도 패키지로 정의되어 있어 다른 패키지에서 사용 시에는 export name으로 정의(대문자 시작) 필요
    ```go
    tracer.initTracing()
    ```
    :arrow_right:
    ```go
    tracer.InitTracing()
    ``` 

- 246 페이지 코드 7-15 `GetProduct()` 메서드 호출 부분 : `wrappers` 사용으로 변경 (proto 및 서비스와 일치)
    ```go
    product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
    ```
    :arrow_right:
    ```go
    product, err := c.GetProduct(ctx, &wrapper.StringValue{Value: r.Value})
    ``` 

---
## 8장 
### 코드 부분
- 257 페이지 맨 아래 `curl` 호출 부분 : 마지막 라인은 명령이 아닌 출력 내용으로 마지막 `\` 부분 삭제
    ```shell
    $ curl -X POST http://localhost:8081/v1/product \
    -d '{"name": "Apple", "description": "iphone7", "price": 699}' \
    "38e13578-d91e-11e9"
    ```
    :arrow_right:
    ```
    $ curl -X POST http://localhost:8081/v1/product \
    -d '{"name": "Apple", "description": "iphone7", "price": 699}'
    ```
    - 마지막 `"38e13578-d91e-11e9"` 라인은 출력된 결과로 생성된 제품 ID입니다.
    
