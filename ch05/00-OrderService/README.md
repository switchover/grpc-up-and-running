# 주문 관리 구현 (Chapter 5 : gRPC: 고급 기능)

5장 고급 기능들은 3장에서 완성된 주문 관리 서비스 및 클라이언트 구현을 활용하여 필요한 기능들을 추가하는 형태로 설명됩니다. 
따라서 아래 기본 구현 및 빌드 등의 과정을 그대로 활용합니다.

## 1. protobuf 정의 파일 생성
[order_management.proto](../order_management.proto)

## 2. Go 서비스용 모듈 생성
Go 모듈을 위한 디렉토리 생성 후, `go mod` 명령을 통해 다음과 같이 모듈을 생성합니다.
```shell
$ mkdir -p order-service/server
$ cd order-service/server
$ go mod init order-service/server
```

## 3. protobuf 파일 복사
별도로 정의된 `order_management.proto` 파일을 `ecommerce` 디렉토리 생성 후 이 디렉토리로 복사합니다.
```shell
$ mkdir ecommerce
$ cp ../../../order_management.proto ecommerce
```
- `order_management.proto`는 임의의 위치에서 복사함 (위 예는 현재 예제 디렉토리 구성의 경우임)

## 4. Go 언어 Skeleton 생성 
다음과 같이 이미 설치된 `protoc` 명령을 통해 skeleton 코드를 생성합니다.
```shell
$ protoc -I ecommerce ecommerce/order_management.proto --go_out=plugins=grpc:ecommerce 
```

## 5. Go 서비스 구현
다음과 같이 Go 서비스를 구현을 참조합니다. 해당 구현 부분은 5장의 고급 기능들을 살펴보기 위한 기본 구현입니다.
[main.go](order-service/server/main.go)

## 6. Go 서버 빌드
다음과 같이 서버를 빌드하고 실행합니다.
```shell
$ go build -i -v -o bin/server main.go
```

## 7. Go 클라이언트 생성
다음과 같인 모듈 생성 및 Stub을 생성합니다.
```shell
$ mkdir -p order-service/client
$ cd order-service/client
$ go mod init order-service/client

$ mkdir ecommerce
$ cp ../../../order_management.proto ecommerce

$ protoc -I ecommerce ecommerce/order_management.proto --go_out=plugins=grpc:ecommerce 
```

## 8. Go 클라이언트 구현 참조
클라이언트 구현은 다음과 같습니다.
[main.go](order-service/client/main.go)

## 9. Go 클라이언트 빌드 및 실행
다음과 같이 클라이언트를 빌드 및 실행합니다.
```shell
$ go build -i -v -o bin/client main.go
$ bin/client
```
