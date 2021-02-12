# Section 4 : 에러 처리 (Chapter 5 : gRPC: 고급 기능)

## 예제 코드 리스트
- 코드 5-7 (서버 에러 처리) : [main.go](order-service/server/main.go)
- 코드 5-8 (클라이언트 에러 처리) : [main.go](order-service/client/main.go)

----
# 에러 처리(Error Handling) (Go)

## 1. 서버 에러 처리(Server error handling) 구현
다음과 같이 Go 서비스에 에러 처리을 적용합니다.
[main.go](order-service/server/main.go) (코드 5-7)

```go
if orderReq.Id == "-1" {
	log.Printf("Order ID is invalid! -> Received Order ID %s",
		orderReq.Id)

	errorStatus := status.New(codes.InvalidArgument,
		"Invalid information received")
	ds, err := errorStatus.WithDetails(
		&epb.BadRequest_FieldViolation{
			Field: "ID",
			Description: fmt.Sprintf(
				"Order ID received is not valid %s : %s",
				orderReq.Id, orderReq.Description),
		},
	)
	if err != nil {
		return nil, errorStatus.Err()
	}

	return nil, ds.Err()
	// ...
```

## 2. 클라이언트 에러 처리(Client error handling) 구현
서버에서의 에러 처리를 확인 및 처리하기 위해 다음과 같이 Go 클라이언트를 수정합니다.
[main.go](order-service/client/main.go) (코드 5-8)

```go
order1 := pb.Order{Id: "-1",
	Items:       []string{"iPhone XS", "Mac Book Pro"},
	Destination: "San Jose, CA", Price: 2300.00}
res, addOrderError := client.AddOrder(ctx, &order1)

if addOrderError != nil {
	errorCode := status.Code(addOrderError)
	if errorCode == codes.InvalidArgument {
		log.Printf("Invalid Argument Error : %s", errorCode)
		errorStatus := status.Convert(addOrderError)
		for _, d := range errorStatus.Details() {
			switch info := d.(type) {
			case *epb.BadRequest_FieldViolation:
				log.Printf("Request Field Invalid: %s", info)
			default:
				log.Printf("Unexpected error type: %s", info)
			}
		}
	} else {
		log.Printf("Unhandled error : %s ", errorCode)
	}
} else {
	log.Print("AddOrder Response -> ", res.Value)
}
```
