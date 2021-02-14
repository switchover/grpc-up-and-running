# 테스팅 (Chapter 7 : 서비스 수준 gRPC 실행)

## 예제 코드 리스트
- 코드 7-1 (서버측 테스트) : [main_test.go](productinfo/server/main_test.go)
- 코드 7-2 (클라이언트 mock 테스트) : [prodinfo_mock_test.go](productinfo/client/mock_prodinfo/prodinfo_mock_test.go)

## 1. 서버측 테스트
gRPC 서비스에 대한 테스트는 다음과 같이 작성됩니다. [main_test.go](productinfo/server/main_test.go) (코드 7-1)
```go
func TestServer_AddProduct(t *testing.T) {
	grpcServer := initGRPCServerHTTP2()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		grpcServer.Stop()
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Sumsung S10"
	description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2019"
	price := float32(700.0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx, &pb.Product{Name: name,
		Description: description, Price: price})
	if err != nil {
		t.Fatalf("Could not add product: %v", err)
	}

	if r.Value == "" {
		t.Errorf("Invalid Product ID %s", r.Value)
	}
	log.Printf("Res %s", r.Value)
	grpcServer.Stop()
}
```

테스트를 실행하기 위해서는 다음과 같이 실행하면 됩니다.

```shell
$ go test
```

## 2. 클라이언트측 테스트
gRPC 클라이언트를 테스트하기 위해서는 Mocking 프레임워크를 사용하는 것이 효율적입니다. 이를 위해 주로 사용하는 [Gomock](https://github.com/golang/mock)을 우선 다음과 같이 설치합니다.

```shell
$ go get github.com/golang/mock/mockgen
```

이제 생성된 Stub 코드에 대하여 다음과 같이 mock 코드를 생성합니다. (모듈 디렉토리로 이동 및 mock을 위한 별도 디렉토리 생성 후 실행)
```shell
$ cd productinfo/client
$ mkdir mock_prodinfo
$ mockgen productinfo/client/ecommerce ProductInfoClient > mock_prodinfo/prodinfo_mock.go
```
참고로 첫번쨰 파라미터는 모듈명(`product/client`)을 포함하여 패키지를 지정하며, 두번째 파리미터는 인터페이스 이름을 지정합니다.

이제 다음과 같이 생성된 mock 코드를 기반으로 테스트를 작성합니다. [prodinfo_mock_test.go](productinfo/client/mock_prodinfo/prodinfo_mock_test.go) (코드 7-2)

```go
func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocklProdInfoClient := NewMockProductInfoClient(ctrl)
	//...
	
	req := &pb.Product{Name: name, Description: description, Price: price}

	mocklProdInfoClient.
		EXPECT().AddProduct(gomock.Any(), &rpcMsg{msg: req}).
		Return(&wrapper.StringValue{Value: "ABC123" + name}, nil)

	testAddProduct(t, mocklProdInfoClient)
}

func testAddProduct(t *testing.T, client pb.ProductInfoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//...
	
	r, err := client.AddProduct(ctx, &pb.Product{Name: name,
		Description: description, Price: price})
	// 테스트와 응답 검증
}
```

테스트를 실행하기 위해서는 다음과 같이 실행하면 됩니다.

```shell
$ go test productinfo/client/mock_prodinfo
```


## 3. 부하 테스트
gRPC 부하 테스트를 위해 [ghz](https://ghz.sh)를 설치합니다.

우선, 다음과 같이 greeter 서비스에 대한 서버 구현을 확인 및 실행합니다. [main.go](loadtest/main.go) 
```shell
$ cd loadtest
$ go build -i -v -o bin/server
$ bin/server
```

그런 다음 proto 파일을 사용해 다음과 같이 부하 테스트를 수행할 수 있습니다.

```shell
$ ghz --insecure \
--proto ./greeter.proto \
--call helloworld.Greeter.SayHello \
-d '{"name":"Joe"}' \
-n 2000 \
-c 20 \
localhost:50051
```
