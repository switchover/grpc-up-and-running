# Chapter 5 : gRPC: 고급 기능

## 예제 코드 리스트
- 코드 5-1 (서버 단일 인터셉터), 코드 5-2 (서버 스트리밍 인터셉터) : [main.go](01-Interceptors/order-service/server/main.go)
- 코드 5-3 (클라이언트 단일 인터셉터), 코드 5-4 (클라이언트 스트리밍 인터셉터) : [main.go](01-Interceptors/order-service/client/main.go)
- 코드 5-5 (데드라인) : [main.go](02-Deadlines/order-service/client/main.go)
- 코드 5-6 (취소 처리) : [main.go](03-Cancellation/order-service/client/main.go)
- 코드 5-7 (서버 에러 처리) : [main.go](04-ErrorHandling/order-service/server/main.go)
- 코드 5-8 (클라이언트 에러 처리) : [main.go](04-ErrorHandling/order-service/client/main.go)
- 코드 5-9 (서버 멀티플렉싱) : [main.go](05-Multiplexing/order-service/server/main.go)
- 코드 5-10 (클라이언트 멀티플렉싱) : [main.go](05-Multiplexing/order-service/client/main.go)
- 코드 5-11 (클라이언트 메타데이터 전송), 코드 5-12 (클라이언트 메타데이터 읽기) : [main.go](06-Metadata/some-service/client/main.go)
- 코드 5-13 (서비스 메타데이터 읽기), 코드 5-14 (서비스 메타데이터 전송) [main.go](06-Metadata/some-service/server/main.go)
- 코드 5-15 (네임 리졸버 구현), 코드 5-16 (클라이언트 로드밸런싱 구현): [main.go](07-LoadBalancing/echo/client/main.go)

## 정오
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
# 세부 세션별 예제

* 기본 주문 서비스 (Go) : [Order Service](./00-OrderService) (5장 전체 공통 사용 예제)
* 인터셉터 (Go) : [Interceptors](./01-Interceptors)
* 데드라인 (Go) : [Deadlines](./02-Deadlines)
* 취소 처리 (Go) : [Cancellation](./03-Cancellation)
* 에러 처리 (Go) : [Error Handling](./04-ErrorHandling)
* 멀티플렉싱 (Go) : [Multiplexing](./05-Multiplexing)
* 메타데이터 (Go) : [Metadata](./06-Metadata)
* 로드밸런싱 (Go) : [Load Balancing](./07-LoadBalancing)
* 압축 (Go) : [Compression](./08-Compression)

---
# 서비스 정의
5장에서는 다음과 같은 proto 파일을 활용합니다. [order_management.proto](order_management.go)
```
syntax = "proto3";

import "google/protobuf/wrappers.proto";

package ecommerce;

service OrderManagement {
    rpc addOrder(Order) returns (google.protobuf.StringValue);
    rpc getOrder(google.protobuf.StringValue) returns (Order);
    rpc searchOrders(google.protobuf.StringValue) returns (stream Order);
    rpc updateOrders(stream Order) returns (google.protobuf.StringValue);
    rpc processOrders(stream google.protobuf.StringValue) returns (stream CombinedShipment);
}

message Order {
    string id = 1;
    repeated string items = 2;
    string description = 3;
    float price = 4;
    string destination = 5;
}

message CombinedShipment {
    string id = 1;
    string status = 2;
    repeated Order ordersList = 3;
}
```

gRPC 고급 기능들을 알아 보기 위해 이미 구현된 주문 관리 서비스를 활용하여, 부가 기능들을 추가하는 형태로 설명됩니다.
주문 서비스 구현은 별도의 [기본 주문 서비스](./00-OrderService)를 참조합니다.

---
# 최종 코드

gRPC의 고급 기능에 대한 예제 코드는 원서의 소스 저장소 [5장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch05)을 참고합니다.
