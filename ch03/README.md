# Chapter 3 : gRPC 통신 패턴

## 예제 코드 리스트
- 코드 3-1 (단순 RPC 서비스 정의 파일) : [order_management.proto](01-SimpleRPC/order_management.proto)
- 코드 3-2 (단순 RPC 서비스 구현) : [main.go](01-SimpleRPC/ordermgt/service/server/main.go)
- 코드 3-3 (단순 RPC 클라이언트 구현) : [main.go](01-SimpleRPC/ordermgt/client/main.go)
- 코드 3-4 (서버 스트리밍 서비스 정의 파일) : [order_management.proto](02-ServerStreamingRPC/order_management.proto)
- 코드 3-5 (서버 스트리밍 서비스 구현) : [main.go](02-ServerStreamingRPC/ordermgt/service/server/main.go)
- 코드 3-6 (서버 스트리밍 클라이언트 구현) : [main.go](02-ServerStreamingRPC/ordermgt/client/main.go)
- 코드 3-7 (클라이언트 스트리밍 서비스 정의 파일) : [order_management.proto](03-ClientStreamingRPC/order_management.proto)
- 코드 3-8 (클라이언트 스트리밍 서비스 구현) : [main.go](03-ClientStreamingRPC/ordermgt/service/server/main.go)
- 코드 3-9 (클라이언트 스트리밍 클라이언트 구현) : [main.go](03-ClientStreamingRPC/ordermgt/client/main.go)
- 코드 3-10 (양방향 스트리밍 서비스 정의 파일) : [order_management.proto](04-BidirectionalStreamingRPC/order_management.proto)
- 코드 3-11 (양방향 스트리밍 서비스 구현) : [main.go](04-BidirectionalStreamingRPC/ordermgt/service/server/main.go)
- 코드 3-12 (양방향 스트리밍 클라이언트 구현) : [main.go](04-BidirectionalStreamingRPC/ordermgt/client/main.go)

## 정오
### 코드 부분
- 98 페이지 코드 3-9. `client` 정의 부분 : 변수명 변경
    ```
    c := pb.NewOrderManagementClient(conn)
    ```
    :arrow_right:
    ```
    client := pb.NewOrderManagementClient(conn)
    ```

---
# 세부 세션별 예제

* 단순 RPC 서비스 구현 (Go) : [Simple RPC](./01-SimpleRPC)
* 서버 스트리밍 RPC 서비스 구현 (Go) : [Server Streaming RPC](./02-ServerStreamingRPC)
* 클라이언트 스트리밍 RPC 서비스 구현 (Go) : [Client Streaming RPC](./03-ClientStreamingRPC)
* 양방향 스트리밍 RPC 서비스 구현 (Go) : [Bidirectional Streaming RPC](./04-BidirectionalStreamingRPC)

---
# 최종 코드

gRPC의 4가지 통신 패턴에 대한 예제 코드는 원서의 소스 저장소 [3장 부분](https://github.com/grpc-up-and-running/samples/tree/master/ch03)을 참고합니다.
