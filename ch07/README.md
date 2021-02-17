# Chapter 7 : 서비스 수준 gRPC 실행

## 예제 코드 리스트
- 코드 7-1 (서버측 테스트) : [main.go](01-Testing/productinfo/server/main_test.go)
- 코드 7-2 (클라이언트 mock 테스트) : [prodinfo_mock_test.go](01-Testing/productinfo/client/mock_prodinfo/prodinfo_mock_test.go)
- 코드 7-4, 7-5 (서버용 쿠버네티스 배포 및 서비스 기술자) : [grpc-prodinfo-server.yaml](03-Kubernetes/productinfo/server/grpc-prodinfo-server.yaml)
- 코드 7-6 (클라이언트용 쿠버네티스 잡 기술자) : [grpc-prodinfo-client-job.yaml](03-Kubernetes/productinfo/client/grpc-prodinfo-client-job.yaml)
- 코드 7-7 (서비스용 쿠버네티스 인그레스 기술자) : [grpc-prodinfo-ingress.yaml](03-Kubernetes/productinfo/ingress/grpc-prodinfo-ingress.yaml)
- 코드 7-8 (서비스 모니터링 활성화) : [main.go](04-OpenCensus/productinfo/server/main.go)
- 코드 7-9 (클라이언트 모니터링 활성화) : [main.go](04-OpenCensus/productinfo/client/main.go)
- 코드 7-10, 7-11 (서비스 모니터링 활성화) : [main.go](05-Prometheus/productinfo/server/main.go)
- 코드 7-12 (클라이언트 모니터링 활성화) : [main.go](05-Prometheus/productinfo/client/main.go)
- 코드 7-13 (예거 익스포터 초기화) : [tracer.go](06-Tracing/productinfo/server/tracer/tracer.go)
- 코드 7-14 (서비스 계측 추가) : [main.go](06-Tracing/productinfo/server/main.go)
- 코드 7-15 (클라아언트 계측 추가) : [main.go](06-Tracing/productinfo/client/main.go)

## 정오
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
# 세부 세션별 예제

* 테스팅 : [Testing](./01-Testing)
* 도커 배포 : [Docker](./02-Docker)
* 쿠버네티스 배포 : [Kubernetes](./03-Kubernetes)
* 오픈센서스 활용 : [OpenCensus](./04-OpenCensus)
* 프로메테우스 활용 : [Prometheus](./05-Prometheus)
* 추적 활용 : [Tracing](./06-Tracing)

---
# 최종 코드

gRPC의 CI/CD, 모니터링 등에 대한 예제 코드는 원서의 소스 저장소 [7장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch07)을 참고합니다.
