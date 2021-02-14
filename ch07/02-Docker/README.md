# 도커 배포 (Chapter 7 : 서비스 수준 gRPC 실행)

## 예제 코드 리스트
- 코드 7-3 (서버용 도커파일) : [Dockerfile](productinfo/server/Dockerfile)

## 1. 서버 도커 배포
다음과 같이 도커 파일을 사용해 도커 컨테이너를 만듭니다. [Dockerfile](productinfo/server/Dockerfile) (코드 7-3)
※ `location` 등의 디렉토리를 각자 시스템에 맞게 수정하셔야 합니다.
```Dockerfile
# 멀티 단계(Multistage) 빌드

# I 단계 빌드:
FROM golang AS build
ENV location /go/src/github.com/grpc-up-and-running/samples/ch07/grpc-docker/go
WORKDIR ${location}/server

ADD ./server ${location}/server
ADD ./proto-gen ${location}/proto-gen

RUN go get -d ./...
RUN go install ./...

RUN CGO_ENABLED=0 go build -o /bin/grpc-productinfo-server

# II 단계 빌드:
FROM scratch
COPY --from=build /bin/grpc-productinfo-server /bin/grpc-productinfo-server

ENTRYPOINT ["/bin/grpc-productinfo-server"]
EXPOSE 50051
```

이제 다음과 같이 도커 명령을 사용해 도커 이미지를 생성합니다.

```shell
$ cd productinfo
$ docker image build -t grpc-productinfo-server -f server/Dockerfile .
```

## 2. 쿨라이언트 도커 배포
클라이언트도 도커 파일을 사용해 도커 컨테이너를 만듭니다. [Dockerfile](productinfo/client/Dockerfile) 

참고로 클라이언트 코드 상에 연결하고자 하는 호스트 정보는 다음과 같이 변경되었습니다.
```go
const (
	address = "productinfo:50051"
)
```
`productinfo`는 docker의 네트워크 호스트명으로 지정됩니다.

그런 다음 다음과 같이 도커 이미지를 생성합니다.

```shell
$ cd productinfo
$ docker image build -t grpc-productinfo-client -f client/Dockerfile .
```

## 3. 도커 실행
우선 클라이언트 컨테이너에서 서버 컨테이너 연결을 위해 네트워크 설정과 함께 다음과 같이 서버 및 클라이언트를 실행합니다.

```shell
$ docker network create my-net

$ docker run -it --network=my-net --name=productinfo --hostname=productinfo -p 50051:50051  grpc-productinfo-server

$ docker run -it --network=my-net --hostname=client grpc-productinfo-client 
```


## 4. 참고
추가로 다음과 같이 docker 이미지를 docker hub 등의 Registry로 등록할 수 있습니다.
```shell
docker image tag grpc-productinfo-server switchover/grpc-productinfo-server
docker image tag grpc-productinfo-client switchover/grpc-productinfo-client
docker image push switchover/grpc-productinfo-server
docker image push switchover/grpc-productinfo-client
```