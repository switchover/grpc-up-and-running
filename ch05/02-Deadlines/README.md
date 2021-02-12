# Section 2 : 데드라인 (Chapter 5 : gRPC: 고급 기능)

## 예제 코드 리스트
- 코드 5-5 (데드라인) : [main.go](order-service/client/main.go)


----
# 데드라인(Deadlines) (Go)

## 1. 클라이언트 데드라인(Client deadlines) 구현
다음과 같이 Go 클라이언트에 데드라인을 적용합니다.
[main.go](order-service/client/main.go) (코드 5-5)

```go
conn, err := grpc.Dial(address, grpc.WithInsecure())
if err != nil {
	log.Fatalf("did not connect: %v", err)
}
defer conn.Close()
client := pb.NewOrderManagementClient(conn)

clientDeadline := time.Now().Add(
	time.Duration(2 * time.Second))
ctx, cancel := context.WithDeadline(
	context.Background(), clientDeadline)

defer cancel()

// Order 등록
order1 := pb.Order{Id: "101",
	Items:       []string{"iPhone XS", "Mac Book Pro"},
	Destination: "San Jose, CA",
	Price:       2300.00}
res, addErr := client.AddOrder(ctx, &order1)
if addErr != nil {
	got := status.Code(addErr)
	log.Printf("Error Occured -> addOrder : , %v:", got)
} else {
	log.Print("AddOrder Response -> ", res.Value)
}
```

## 2. 클라이언트 데이터라인(Client deadlines)용 서버 구현 수정
데드라인 테스트를 위해 서버에 3초의 sleep을 추가하여 데드라인을 강제로 발생시킵니다.
[main.go](order-service/server/main.go)

```go
// Order 등록
func (s *server) AddOrder(ctx context.Context, orderReq *pb.Order) (*wrappers.StringValue, error) {
	//----------------------------------
	// 데드라인 테스트용 Sleep 추가
	//----------------------------------
	log.Printf("Sleep 3 seconds...")
	time.Sleep(time.Second * 3)
	//----------------------------------
	orderMap[orderReq.Id] = *orderReq
	log.Printf("Order : %v -> Added", orderReq.Id)
	return &wrappers.StringValue{Value: "Order Added: " + orderReq.Id}, nil
}
```
