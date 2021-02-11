# Chapter 5 : gRPC: 고급 기능

## 예제 코드 리스트
- 코드 5-1 (서버 단일 인터셉터), 코드 5-2 (서버 스트리밍 인터셉터) : [main.go](01-Interceptors/order-service/server/main.go)
- 코드 5-3 (클라이언트 단일 인터셉터), 코드 5-4 (클라이언트 스트리밍 인터셉터) : [main.go](01-Interceptors/order-service/client/main.go)

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
# 세부 세션별 예제

* 기본 주문 서비스 (Go) : [Order Service](./00-OrderService)
* 인터셉터 (Go) : [Interceptors](./01-Interceptors)

---
# 최종 코드

gRPC의 고급 기능에 대한 예제 코드는 원서의 소스 저장소 [5장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch05)을 참고합니다.
