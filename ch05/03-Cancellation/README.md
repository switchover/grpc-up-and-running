# Section 3 : 취소 처리 (Chapter 5 : gRPC: 고급 기능)

## 예제 코드 리스트
- 코드 5-6 (취소 처리) : [main.go](order-service/client/main.go)


----
# 취소 처리(Cancellation) (Go)

## 1. 클라이언트 취소 처리(Client cancellation) 구현
다음과 같이 Go 클라이언트에 취소 처리를 적용합니다.
[main.go](order-service/client/main.go) (코드 5-6)

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
defer cancel()

streamProcOrder, _ := client.ProcessOrders(ctx)
_ = streamProcOrder.Send(&wrapper.StringValue{Value: "102"})
_ = streamProcOrder.Send(&wrapper.StringValue{Value: "103"})
_ = streamProcOrder.Send(&wrapper.StringValue{Value: "104"})

channel := make(chan bool, 1)
go asncClientBidirectionalRPC(streamProcOrder, channel)
time.Sleep(time.Millisecond * 1000)

// RPC 취소
cancel()
log.Printf("RPC Status : %s", ctx.Err())

_ = streamProcOrder.Send(&wrapper.StringValue{Value: "101"})
_ = streamProcOrder.CloseSend()

<-channel

// ...
func asncClientBidirectionalRPC (
	streamProcOrder pb.OrderManagement_ProcessOrdersClient, c chan bool) {
	// ...
	combinedShipment, errProcOrder := streamProcOrder.Recv()
	if errProcOrder != nil {
		log.Printf("Error Receiving messages %v", errProcOrder)
	// ...
}
```

## 2. 클라이언트 취소 처리(Client cancellation)용 서버 구현 수정
취소 처리를 위한 서버의 수정 사항은 없습니다.
